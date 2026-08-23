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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
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
