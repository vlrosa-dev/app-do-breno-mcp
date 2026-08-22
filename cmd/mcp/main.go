package main

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/clients"
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/tools"
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
	)
	token, err := api.Login(
		context.Background(),
		config.APIUsername,
		config.APIPassword,
	)
	api.SetToken(token.AccessToken)
	userInfo, err := api.FetchUserInfo(
		context.Background(),
	)
	api.SetUserId(userInfo.ID)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "app-do-breno-mcp",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_transaction",
		Description: "Create a new transaction",
	}, tools.ToolCreateTransaction(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch_transactions",
		Description: "Fetch all transactions",
	}, tools.ToolFetchTransactions(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_transaction",
		Description: "Delete a transaction",
	}, tools.ToolDeleteTransaction(api))

	// Run server on stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
