package main

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/clients"
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/structs"
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
		config.APIUsername,
		config.APIPassword,
	)
	// id, _ := api.CreateExpenseOne(
	// 	context.Background(),
	// 	&structs.DespesasUnicaRequest{
	// 		Data:      "2026-08-23",
	// 		Descricao: "Teste 3",
	// 		TagIds:    []int{},
	// 		Tipo:      "Diario",
	// 		User:      api.UserID,
	// 		Valor:     1.50,
	// 	},
	// )
	// log.Println("ID: ", id)

	err = api.DeleteExpenseOne(
		context.Background(),
		&structs.DespesasUnicaDeleteRequest{
			Id:     "966c852d-fe45-4ca4-a623-876cff8a2db5",
			UserId: api.UserID,
		},
	)
	if err != nil {
		log.Println("err: ", err.Error())
	}

	// list, _ := api.ListExpenses(
	// 	context.Background(),
	// 	structs.DespesasListarRequest{
	// 		Select:      "*",
	// 		UserId:      userInfo.ID,
	// 		Tipo:        "Diario",
	// 		DataInicial: "2026-08-01",
	// 		DataFinal:   "2026-09-01",
	// 	},
	// )
	// log.Println(list)

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
