package api

import (
	"app-do-breno-agent/internal/structs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (api *API) GetToken(ctx context.Context, email, password string) (structs.AuthResponse, error) {
	fullUrl := fmt.Sprintf("%s%s", api.baseURL, "/auth/v1/token")

	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.AuthResponse{}, err
	}

	queryParams := url.Values{}
	queryParams.Add("grant_type", "password")
	u.RawQuery = queryParams.Encode()

	structPayload := structs.AuthRequest{
		Email:    email,
		Password: password,
	}

	jsonPayload, err := json.Marshal(structPayload)
	if err != nil {
		return structs.AuthResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(jsonPayload)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.AuthResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+api.apiKey)
	req.Header.Add("ApiKey", api.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		return structs.AuthResponse{}, err
	}

	if res == nil {
		return structs.AuthResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.AuthResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return structs.AuthResponse{}, msgError
	}

	var authResponse structs.AuthResponse
	err = json.Unmarshal(responseData, &authResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.AuthResponse{}, err
	}

	return authResponse, nil
}
