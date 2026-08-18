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

func (api *API) CriarMovimentacao(
	ctx context.Context, apiKey, token string, movimentacao *structs.MovimentacaoRequest,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url, err := url.Parse(api.BaseURL + "/rest/v1/rpc/movimentacao_create_with_tags_v2")
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Apikey", apiKey)

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
