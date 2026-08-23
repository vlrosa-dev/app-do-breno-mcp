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
	token := auth.AccessToken
	log.Println("Token: ", token[:30])

	user, err := GetUser(
		context.Background(),
		client,
		baseURL,
		token,
		apiKey,
	)
	if err != nil {
		log.Fatalf("erro get user info, msg: %v", err.Error())
	}
	userId := user.ID
	log.Println("UserID: ", userId)

	return &API{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
		token:   token,
		userId:  userId,
	}
}

func (api *API) SetToken(token string) {
	api.token = token
}

func (api *API) SetUserId(userId string) {
	api.userId = userId
}
