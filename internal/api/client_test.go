package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type fakeAuth struct {
	token      string
	next       string
	forceCalls int32
	tokenErr   error
}

func (f *fakeAuth) Token(context.Context) (string, error) {
	if f.tokenErr != nil {
		return "", f.tokenErr
	}
	return f.token, nil
}

func (f *fakeAuth) ForceLogin(context.Context) (string, error) {
	atomic.AddInt32(&f.forceCalls, 1)
	f.token = f.next
	return f.next, nil
}

func newTestClient(t *testing.T, h http.HandlerFunc, auth TokenProvider) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c := NewClient(ClientOptions{
		BaseURL:   srv.URL,
		APIKey:    "test-api-key",
		Auth:      auth,
		Transport: srv.Client().Transport,
	})
	return c, srv
}

func TestREST_GetDecodesAndSendsAuthHeaders(t *testing.T) {
	var gotAuth, gotKey, gotQuery, gotMethod, gotPath string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotKey = r.Header.Get("apikey")
		gotQuery = r.URL.RawQuery
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"nome":"Padaria"}`)
	}, &fakeAuth{token: "tok-1"})

	var out struct {
		Nome string `json:"nome"`
	}
	q := map[string][]string{"select": {"*"}}
	if err := c.REST(context.Background(), http.MethodGet, "/rest/v1/cartao", q, nil, &out); err != nil {
		t.Fatalf("REST: %v", err)
	}
	if gotMethod != http.MethodGet || gotPath != "/rest/v1/cartao" {
		t.Fatalf("método/path = %s %s", gotMethod, gotPath)
	}
	if gotAuth != "Bearer tok-1" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotKey != "test-api-key" {
		t.Fatalf("apikey = %q", gotKey)
	}
	if gotQuery != "select=%2A" {
		t.Fatalf("query = %q", gotQuery)
	}
	if out.Nome != "Padaria" {
		t.Fatalf("out.Nome = %q", out.Nome)
	}
}

func TestRPC_PostsJSONBody(t *testing.T) {
	var gotPath, gotCT string
	var gotBody []byte
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotCT = r.Header.Get("Content-Type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `"ok"`)
	}, &fakeAuth{token: "t"})

	in := map[string]any{"p_valor": 1.5, "p_descricao": "Padaria"}
	var out string
	if err := c.RPC(context.Background(), "movimentacao_create_with_tags_v2", in, &out); err != nil {
		t.Fatalf("RPC: %v", err)
	}
	if gotPath != "/rest/v1/rpc/movimentacao_create_with_tags_v2" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotCT != "application/json" {
		t.Fatalf("content-type = %q", gotCT)
	}
	if string(gotBody) != `{"p_descricao":"Padaria","p_valor":1.5}` {
		t.Fatalf("body = %s", gotBody)
	}
	if out != "ok" {
		t.Fatalf("out = %q", out)
	}
}

func TestREST_NoContent(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}, &fakeAuth{token: "t"})

	out := "untouched"
	if err := c.REST(context.Background(), http.MethodDelete, "/rest/v1/movimentacao", nil, nil, &out); err != nil {
		t.Fatalf("REST: %v", err)
	}
	if out != "untouched" {
		t.Fatalf("out mutado: %q", out)
	}
}

func TestREST_ErrorStatusReturnsAPIError(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"message":"tipo invalido"}`)
	}, &fakeAuth{token: "t"})

	err := c.REST(context.Background(), http.MethodPost, "/rest/v1/rpc/x", map[string][]string{}, map[string]any{}, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("esperava *APIError, veio %T: %v", err, err)
	}
	if apiErr.Status != http.StatusBadRequest {
		t.Fatalf("status = %d", apiErr.Status)
	}
	if apiErr.Body != `{"message":"tipo invalido"}` {
		t.Fatalf("body = %q", apiErr.Body)
	}
}

func TestREST_RetriesOnceAfter401(t *testing.T) {
	var calls int32
	auth := &fakeAuth{token: "stale", next: "fresh"}
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		body, _ := io.ReadAll(r.Body)
		if r.Header.Get("Authorization") != "Bearer fresh" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if string(body) != `{"p_valor":2}` {
			t.Errorf("corpo não rebobinado no retry: %s", body)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `"done"`)
	}, auth)

	var out string
	if err := c.RPC(context.Background(), "fn", map[string]any{"p_valor": 2}, &out); err != nil {
		t.Fatalf("RPC: %v", err)
	}
	if calls != 2 {
		t.Fatalf("esperava 2 chamadas ao servidor, veio %d", calls)
	}
	if atomic.LoadInt32(&auth.forceCalls) != 1 {
		t.Fatalf("esperava 1 ForceLogin, veio %d", auth.forceCalls)
	}
	if out != "done" {
		t.Fatalf("out = %q", out)
	}
}

func TestREST_Persistent401ReturnsAPIError(t *testing.T) {
	auth := &fakeAuth{token: "a", next: "b"}
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `unauthorized`)
	}, auth)

	err := c.RPC(context.Background(), "fn", map[string]any{"x": 1}, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("esperava *APIError, veio %T: %v", err, err)
	}
	if apiErr.Status != http.StatusUnauthorized {
		t.Fatalf("status = %d", apiErr.Status)
	}
	if atomic.LoadInt32(&auth.forceCalls) != 1 {
		t.Fatalf("esperava exatamente 1 ForceLogin, veio %d", auth.forceCalls)
	}
}

func TestREST_TokenErrorPropagates(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("servidor não deveria ser chamado")
	}, &fakeAuth{tokenErr: errors.New("sem credenciais")})

	if err := c.Get(context.Background(), "/rest/v1/cartao", nil, nil); err == nil {
		t.Fatal("esperava erro de token")
	}
}
