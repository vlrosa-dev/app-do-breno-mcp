package api

import (
	"context"
	"log"
	"net/http"
)

type API struct {
	BaseUrl  string
	Client   *http.Client
	ApiToken string
	ApiKey   string
	UserID   string
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
		BaseUrl:  baseURL,
		ApiKey:   apiKey,
		Client:   client,
		ApiToken: token,
		UserID:   userId,
	}
}
