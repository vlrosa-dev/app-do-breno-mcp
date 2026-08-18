package www

import (
	"app-do-breno-agent/internal/structs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	baseURL string
	client  *http.Client
	token   string
	apiKey  string
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

func (api *API) Login(ctx context.Context, email, password string) (*structs.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.baseURL + "/auth/v1/token")
	if err != nil {
		return nil, fmt.Errorf("parsing base URL: %w", err)
	}

	queryParams := url.Query()
	queryParams.Set("grant_type", "password")
	url.RawQuery = queryParams.Encode()

	structPayload := structs.AuthRequest{
		Email:    email,
		Password: password,
	}

	jsonPayload, err := json.Marshal(structPayload)
	if err != nil {
		return nil, fmt.Errorf("serializando payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("montando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.apiKey))
	req.Header.Set("Apikey", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chamando api de autenticação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		corpoErro, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // limita leitura de corpo de erro
		return nil, fmt.Errorf("api retornou %d: %s", resp.StatusCode, corpoErro)
	}

	var auth structs.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("decodificando resposta: %w", err)
	}
	return &auth, nil
}
