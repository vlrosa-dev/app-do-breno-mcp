package api

import (
	"context"
	"log"
	"net/http"
)

type API struct {
	baseURL string
	client  *http.Client
	token   string
	apiKey  string
	userId  string
}

func NewAPI(baseURL, apiKey string, client *http.Client, email, password string) *API {
	auth, err := GetToken(
		context.Background(),
		client,
		baseURL,
		apiKey,
		email,
		password,
	)
	if err != nil {
		log.Fatalf("erro authenticate client API, msg: %v", err.Error())
	}

	return &API{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
		token:   auth.AccessToken,
	}
}

func (api *API) SetToken(token string) {
	api.token = token
}

func (api *API) SetUserId(userId string) {
	api.userId = userId
}
