package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// gotrue.go concentra as chamadas cruas aos endpoints /auth/v1/* do Supabase
// (GoTrue). Usa um *http.Client próprio — nunca o client com auto-refresh — para
// evitar recursão durante o login.

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (r authResponse) expiry() time.Time {
	if r.ExpiresAt > 0 {
		return time.Unix(r.ExpiresAt, 0)
	}
	if r.ExpiresIn > 0 {
		return time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)
	}
	return time.Time{}
}

// passwordLogin executa POST /auth/v1/token?grant_type=password.
func passwordLogin(ctx context.Context, hc *http.Client, baseURL, apiKey string, c Credentials) (authResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	u := fmt.Sprintf("%s/auth/v1/token?%s", baseURL, url.Values{"grant_type": {"password"}}.Encode())

	payload, err := json.Marshal(c)
	if err != nil {
		return authResponse{}, fmt.Errorf("serializando credenciais: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(payload))
	if err != nil {
		return authResponse{}, err
	}
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := hc.Do(req)
	if err != nil {
		return authResponse{}, err
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return authResponse{}, fmt.Errorf("login falhou (%s): %s", res.Status, body)
	}

	var out authResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return authResponse{}, fmt.Errorf("parseando resposta de login: %w", err)
	}
	if out.AccessToken == "" {
		return authResponse{}, fmt.Errorf("login sem access_token na resposta")
	}
	return out, nil
}

// fetchUserID executa GET /auth/v1/user e devolve o id do usuário autenticado.
func fetchUserID(ctx context.Context, hc *http.Client, baseURL, apiKey, accessToken string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/auth/v1/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("consulta de usuário falhou (%s): %s", res.Status, body)
	}

	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("parseando resposta de usuário: %w", err)
	}
	if out.ID == "" {
		return "", fmt.Errorf("resposta de usuário sem id")
	}
	return out.ID, nil
}
