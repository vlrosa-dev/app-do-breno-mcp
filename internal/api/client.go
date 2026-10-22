package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vlrosa-dev/app-do-breno-mcp/internal/clients"
)

// TokenProvider entrega um bearer token válido para o Client. É satisfeito por
// *auth.Manager.
type TokenProvider interface {
	// Token devolve um access token válido, reautenticando se necessário.
	Token(ctx context.Context) (string, error)
	// ForceLogin ignora qualquer cache e reautentica. Chamado após um 401.
	ForceLogin(ctx context.Context) (string, error)
}

// Client abstrai as requisições HTTP contra a API (Supabase/PostgREST). Injeta
// `apikey` + `Authorization: Bearer` automaticamente e reautentica de forma
// transparente quando a API responde 401. Inspirado em rhysd/gh-cli/api/client.go.
type Client struct {
	http    *http.Client
	baseURL string
}

// ClientOptions configura o Client.
type ClientOptions struct {
	BaseURL string
	APIKey  string
	Auth    TokenProvider
	// Transport base opcional (útil em testes). Se nil, usa clients.NewTransport().
	Transport http.RoundTripper
	// Timeout total por requisição. Se zero, usa 30s.
	Timeout time.Duration
}

// NewClient monta um Client. `Auth` e `APIKey` são obrigatórios.
func NewClient(opts ClientOptions) *Client {
	base := opts.Transport
	if base == nil {
		base = clients.NewTransport()
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		http: &http.Client{
			Transport: &authTransport{base: base, auth: opts.Auth, apiKey: opts.APIKey},
			Timeout:   timeout,
		},
		baseURL: strings.TrimRight(opts.BaseURL, "/"),
	}
}

// APIError representa uma resposta HTTP com status >= 300.
type APIError struct {
	Method string
	URL    string
	Status int
	Body   string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("%s %s: %d", e.Method, e.URL, e.Status)
	}
	return fmt.Sprintf("%s %s: %d — %s", e.Method, e.URL, e.Status, e.Body)
}

// REST executa uma requisição e, quando `out != nil` e há corpo, faz o unmarshal
// da resposta. `body` é serializado como JSON (nil = sem corpo).
func (c *Client) REST(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var payload []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("serializando corpo de %s %s: %w", method, path, err)
		}
		payload = b
	}

	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reader io.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("lendo resposta de %s %s: %w", method, u, err)
	}

	if res.StatusCode >= http.StatusMultipleChoices {
		return &APIError{Method: method, URL: u, Status: res.StatusCode, Body: strings.TrimSpace(string(data))}
	}

	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parseando resposta de %s %s: %w", method, u, err)
	}
	return nil
}

// RPC chama uma função PostgREST (POST /rest/v1/rpc/<fn>).
func (c *Client) RPC(ctx context.Context, fn string, body, out any) error {
	return c.REST(ctx, http.MethodPost, "/rest/v1/rpc/"+fn, nil, body, out)
}

// Get é um atalho para REST GET.
func (c *Client) Get(ctx context.Context, path string, query url.Values, out any) error {
	return c.REST(ctx, http.MethodGet, path, query, nil, out)
}
