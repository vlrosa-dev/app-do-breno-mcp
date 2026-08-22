package api

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

func (api *API) CreateNewTransaction(
	ctx context.Context, movimentacao *structs.DespesasUnicaRequest,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.baseURL + "/rest/v1/rpc/movimentacao_create_with_tags_v2")
	if err != nil {
		return "", fmt.Errorf("parsing base URL: %w", err)
	}

	jsonPayload, err := json.Marshal(movimentacao)
	if err != nil {
		return "", fmt.Errorf("serializando payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewReader(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("montando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))
	req.Header.Set("Apikey", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("chamando api de autenticação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		corpoErro, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // limita leitura de corpo de erro
		return "", fmt.Errorf("api retornou %d: %s", resp.StatusCode, corpoErro)
	}

	id, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // limita leitura de corpo de erro
	return string(id), nil
}

func (api *API) FetchTransactions(
	ctx context.Context, params *structs.DespesasListarRequest,
) ([]structs.DespesasListarResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.baseURL + "/rest/v1/movimentacao")
	if err != nil {
		return []structs.DespesasListarResponse{}, fmt.Errorf("parsing base URL: %w", err)
	}

	query := url.Query()
	query.Set("select", "*")
	query.Set("user_id", fmt.Sprintf("eq.%s", api.userId))
	query.Set("tipo", fmt.Sprintf("eq.%s", params.Tipo))
	query.Set("data", fmt.Sprintf("gte.%s", params.DataInicial))
	query.Set("data", fmt.Sprintf("lt.%s", params.DataFinal))
	query.Set("order", "data.desc,created_at.desc")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return []structs.DespesasListarResponse{}, fmt.Errorf("montando request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))
	req.Header.Set("Apikey", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return []structs.DespesasListarResponse{}, fmt.Errorf("chamando api de autenticação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		corpoErro, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // limita leitura de corpo de erro
		return []structs.DespesasListarResponse{}, fmt.Errorf("api retornou %d: %s", resp.StatusCode, corpoErro)
	}

	var movimentacoes []structs.DespesasListarResponse
	err = json.NewDecoder(resp.Body).Decode(&movimentacoes)
	if err != nil {
		return []structs.DespesasListarResponse{}, fmt.Errorf("decodificando resposta: %w", err)
	}

	return movimentacoes, nil
}

func (api *API) DeleteTransaction(ctx context.Context, movimentacaoId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.baseURL + "/rest/v1/movimentacao")
	if err != nil {
		return fmt.Errorf("parsing base URL: %w", err)
	}

	query := url.Query()
	query.Set("movimentacao_id", fmt.Sprintf("eq.%s", movimentacaoId))
	query.Set("user_id", fmt.Sprintf("eq.%s", api.userId))
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return fmt.Errorf("montando request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))
	req.Header.Set("Apikey", api.apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		corpoErro, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) // limita leitura de corpo de erro
		return fmt.Errorf("api retornou %d: %s", resp.StatusCode, corpoErro)
	}

	return nil
}
