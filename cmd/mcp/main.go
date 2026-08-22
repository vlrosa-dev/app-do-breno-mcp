package main

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/clients"
	"app-do-breno-agent/internal/config"
	"context"
	"log"

	"github.com/joho/godotenv"
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
	token, err := api.GetToken(
		context.Background(),
		config.APIUsername,
		config.APIPassword,
	)
	log.Println("apiToken:", token.AccessToken)
	api.SetToken(token.AccessToken)
	userInfo, err := api.GetUserInfo(
		context.Background(),
	)
	log.Println("userID:", userInfo.ID)
	api.SetUserId(userInfo.ID)

	// server := mcp.NewServer(&mcp.Implementation{
	// 	Name:    "app-do-breno-mcp",
	// 	Version: "1.0.0",
	// }, nil)

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "create_transaction",
	// 	Description: "Create a new transaction",
	// }, tools.ToolCreateTransaction(api))

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "fetch_transactions",
	// 	Description: "Fetch all transactions",
	// }, tools.ToolFetchTransactions(api))

	// mcp.AddTool(server, &mcp.Tool{
	// 	Name:        "delete_transaction",
	// 	Description: "Delete a transaction",
	// }, tools.ToolDeleteTransaction(api))

	// // Run server on stdio transport
	// if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
	// 	log.Fatal(err)
	// }
}
