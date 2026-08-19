package main

import (
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/tools"
	"app-do-breno-agent/internal/www"
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

	apiClient := www.NewApiClient()

	api := www.NewAPI(
		config.APIBaseURL,
		config.APIKey,
		apiClient,
	)
	token, err := api.Login(
		context.Background(),
		config.APIUsername,
		config.APIPassword,
	)
	api.SetToken(token.AccessToken)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "app-do-breno-mcp",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_transaction",
		Description: "Create a new transaction",
	}, tools.ToolMcpMovimentacao(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_transactions",
		Description: "List all transactions",
	}, tools.ToolMcpListarMovimentacoes(api))

	// Run server on stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
