package api

import (
	"app-do-breno-agent/internal/helpers"
	"app-do-breno-agent/internal/structs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (api *API) GetToken(ctx context.Context, email, password string) (structs.AuthResponse, error) {
	fullUrl := fmt.Sprintf("%s%s", api.baseURL, "/auth/v1/token")

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + api.apiKey,
		"Apikey":        api.apiKey,
	}

	queryParams := url.Values{}
	queryParams.Add("grant_type", "password")

	structPayload := structs.AuthRequest{
		Email:    email,
		Password: password,
	}

	jsonPayload, err := json.Marshal(structPayload)
	if err != nil {
		return structs.AuthResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(jsonPayload)

	var authResponse structs.AuthResponse
	resp, err := helpers.MakeHTTPRequest(
		api.client,
		context.Background(),
		fullUrl,
		http.MethodPost,
		headers,
		queryParams,
		body,
		authResponse,
	)
	if err != nil {
		return structs.AuthResponse{}, fmt.Errorf("request error, msg: %v", err.Error())
	}

	return resp, nil
}
