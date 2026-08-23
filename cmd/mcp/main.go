package main

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/clients"
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/mcpserver"
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
		api,
	)
	if err != nil {
		log.Fatalf("error creating MCP Server, error: %v", err.Error())
	}

	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
