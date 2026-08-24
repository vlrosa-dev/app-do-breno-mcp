package main

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/clients"
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/mcpserver"
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// stdout is reserved for the MCP stdio protocol stream, so all logging must go to stderr.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// .env is optional: when absent, env vars are expected to come from the MCP client config.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logger.Warn("error loading .env file", "error", err.Error())
	}
	config := config.Load()

	apiClient := clients.NewApiClient()

	api := api.NewAPI(
		config.APIBaseURL,
		config.APIKey,
		apiClient,
		config.APIUsername,
		config.APIPassword,
	)

	mcpServer, err := mcpserver.NewMcpServer(
		logger,
		api,
	)
	if err != nil {
		logger.Error("error creating MCP Server", "error: ", err.Error())
	}

	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Error(err.Error())
	}
}
