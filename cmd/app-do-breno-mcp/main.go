package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/vlrosa-dev/app-do-breno-mcp/internal/api"
	"github.com/vlrosa-dev/app-do-breno-mcp/internal/auth"
	"github.com/vlrosa-dev/app-do-breno-mcp/internal/breno"
	"github.com/vlrosa-dev/app-do-breno-mcp/internal/config"
	"github.com/vlrosa-dev/app-do-breno-mcp/internal/mcpserver"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// stdout é reservado para o stream do protocolo MCP; todo log vai para stderr.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// .env é opcional: quando ausente, as variáveis vêm da config do cliente MCP.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logger.Warn("erro carregando .env", "error", err.Error())
	}
	cfg := config.Load()

	storePath, err := auth.DefaultPath()
	if err != nil {
		logger.Error("resolvendo caminho de config", "error", err.Error())
		os.Exit(1)
	}

	manager, err := auth.NewManager(auth.Options{
		BaseURL:     cfg.APIBaseURL,
		APIKey:      cfg.APIKey,
		Store:       auth.NewFileStore(storePath),
		Credentials: auth.Credentials{Email: cfg.Email, Password: cfg.Password},
	})
	if err != nil {
		logger.Error("inicializando autenticação", "error", err.Error())
		os.Exit(1)
	}

	client := api.NewClient(api.ClientOptions{
		BaseURL: cfg.APIBaseURL,
		APIKey:  cfg.APIKey,
		Auth:    manager,
	})

	svc := breno.NewService(client, manager)

	mcpServer, err := mcpserver.NewMcpServer(logger, svc)
	if err != nil {
		logger.Error("erro criando o servidor MCP", "error", err.Error())
		os.Exit(1)
	}

	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Error(err.Error())
	}
}
