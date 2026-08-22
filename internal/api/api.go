package api

import (
	"net/http"
)

type API struct {
	baseURL string
	client  *http.Client
	token   string
	apiKey  string
	userId  string
}

func NewAPI(baseURL, apiKey string, client *http.Client) *API {
	return &API{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
	}
}

func (api *API) SetToken(token string) {
	api.token = token
}

func (api *API) SetUserId(userId string) {
	api.userId = userId
}
