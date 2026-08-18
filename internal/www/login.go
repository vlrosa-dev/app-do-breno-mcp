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
	BaseURL string
	client  *http.Client
}

func NewAPI(baseURL string, client *http.Client) *API {
	return &API{
		BaseURL: baseURL,
		client:  client,
	}
}

func (api *API) GetToken(ctx context.Context, email, password, apiKey string) (*structs.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.BaseURL + "/auth/v1/token")
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Apikey", apiKey)

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
