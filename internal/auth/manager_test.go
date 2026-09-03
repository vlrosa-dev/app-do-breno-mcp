package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

// authServer simula os endpoints /auth/v1/token e /auth/v1/user do GoTrue.
type authServer struct {
	*httptest.Server
	tokenHits int32
	userHits  int32
	issue     func() authResponse
}

func newAuthServer(t *testing.T) *authServer {
	t.Helper()
	as := &authServer{}
	as.issue = func() authResponse {
		return authResponse{
			AccessToken: "access-" + time.Now().Format("150405.000000000"),
			ExpiresAt:   time.Now().Add(time.Hour).Unix(),
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/v1/token", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&as.tokenHits, 1)
		if r.URL.Query().Get("grant_type") != "password" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var c Credentials
		_ = json.NewDecoder(r.Body).Decode(&c)
		if c.Email == "" || c.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(as.issue())
	})
	mux.HandleFunc("/auth/v1/user", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&as.userHits, 1)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "user-123"})
	})
	as.Server = httptest.NewServer(mux)
	t.Cleanup(as.Close)
	return as
}

func newManager(t *testing.T, as *authServer, store Store, creds Credentials) *Manager {
	t.Helper()
	m, err := NewManager(Options{
		HTTPClient:  as.Client(),
		BaseURL:     as.URL,
		APIKey:      "anon-key",
		Store:       store,
		Credentials: creds,
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestManager_TokenLoginsAndPersists(t *testing.T) {
	as := newAuthServer(t)
	store := NewFileStore(filepath.Join(t.TempDir(), "config.json"))
	m := newManager(t, as, store, Credentials{Email: "a@b.com", Password: "pw"})

	tok, err := m.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if tok == "" {
		t.Fatal("token vazio")
	}
	if m.UserID() != "user-123" {
		t.Fatalf("UserID = %q", m.UserID())
	}

	f, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if f.Email != "a@b.com" || f.Password != "pw" {
		t.Fatalf("credenciais não persistidas: %+v", f)
	}
	if f.Token == nil || f.Token.AccessToken != tok {
		t.Fatalf("token não persistido: %+v", f.Token)
	}
}

func TestManager_TokenUsesCacheWhenValid(t *testing.T) {
	as := newAuthServer(t)
	store := NewFileStore(filepath.Join(t.TempDir(), "config.json"))
	m := newManager(t, as, store, Credentials{Email: "a@b.com", Password: "pw"})

	first, err := m.Token(context.Background())
	if err != nil {
		t.Fatalf("Token 1: %v", err)
	}
	second, err := m.Token(context.Background())
	if err != nil {
		t.Fatalf("Token 2: %v", err)
	}
	if first != second {
		t.Fatalf("token mudou apesar de válido: %q != %q", first, second)
	}
	if hits := atomic.LoadInt32(&as.tokenHits); hits != 1 {
		t.Fatalf("esperava 1 login, veio %d", hits)
	}
}

func TestManager_TokenReloginsWhenExpired(t *testing.T) {
	as := newAuthServer(t)
	store := NewFileStore(filepath.Join(t.TempDir(), "config.json"))

	// Semeia o store com um token já expirado + credenciais.
	seed := File{
		Email:    "a@b.com",
		Password: "pw",
		Token:    &Token{AccessToken: "velho", ExpiresAt: time.Now().Add(-time.Minute), UserID: "user-123"},
	}
	if err := store.Save(seed); err != nil {
		t.Fatalf("seed: %v", err)
	}

	m := newManager(t, as, store, Credentials{})
	tok, err := m.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if tok == "velho" {
		t.Fatal("deveria ter reautenticado")
	}
	if hits := atomic.LoadInt32(&as.tokenHits); hits != 1 {
		t.Fatalf("esperava 1 relogin, veio %d", hits)
	}
}

func TestManager_TokenWithoutCredentials(t *testing.T) {
	as := newAuthServer(t)
	store := NewFileStore(filepath.Join(t.TempDir(), "config.json"))
	m := newManager(t, as, store, Credentials{})

	if _, err := m.Token(context.Background()); err == nil {
		t.Fatal("esperava erro por falta de credenciais")
	}
}

func TestManager_ForceLoginAlwaysHitsServer(t *testing.T) {
	as := newAuthServer(t)
	store := NewFileStore(filepath.Join(t.TempDir(), "config.json"))
	m := newManager(t, as, store, Credentials{Email: "a@b.com", Password: "pw"})

	if _, err := m.Token(context.Background()); err != nil {
		t.Fatalf("Token: %v", err)
	}
	if _, err := m.ForceLogin(context.Background()); err != nil {
		t.Fatalf("ForceLogin: %v", err)
	}
	if hits := atomic.LoadInt32(&as.tokenHits); hits != 2 {
		t.Fatalf("esperava 2 logins, veio %d", hits)
	}
}

func TestManager_LogoutClearsStore(t *testing.T) {
	as := newAuthServer(t)
	path := filepath.Join(t.TempDir(), "config.json")
	store := NewFileStore(path)
	m := newManager(t, as, store, Credentials{Email: "a@b.com", Password: "pw"})

	if _, err := m.Token(context.Background()); err != nil {
		t.Fatalf("Token: %v", err)
	}
	if err := m.Logout(); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, ok := m.Status(); ok {
		t.Fatal("status ainda válido após logout")
	}
	f, _ := store.Load()
	if f.Token != nil || f.Email != "" {
		t.Fatalf("store não limpo: %+v", f)
	}
}
