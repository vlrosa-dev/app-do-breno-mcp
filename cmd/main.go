package main

import (
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/structs"
	"app-do-breno-agent/internal/www"
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
	log.Println(config)

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

	if err != nil {
		log.Fatal("Error fetching token:", err)
	}
	log.Println("Fetched token:", token)

	id_transacao, err := api.CriarMovimentacao(
		context.Background(),
		&structs.MovimentacaoRequest{
			Data:      "2026-08-18",
			Descricao: "Teste de movimentação",
			TagIds:    []int{},
			Tipo:      "Saida",
			User:      "ab1d14db-b3a5-4951-a990-d6e4ac3195a2",
			Valor:     0.02,
		},
	)
	if err != nil {
		log.Fatal("Error creating transaction:", err)
	}
	log.Println("ID da transação criada:", id_transacao)

}
