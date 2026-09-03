package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileStore_LoadMissingReturnsEmpty(t *testing.T) {
	s := NewFileStore(filepath.Join(t.TempDir(), "nao-existe", "config.json"))
	f, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if f.Email != "" || f.Token != nil {
		t.Fatalf("esperava File vazio, veio %+v", f)
	}
}

func TestFileStore_SaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".breno", "config.json")
	s := NewFileStore(path)

	want := File{
		Email:    "user@exemplo.com",
		Password: "segredo",
		BaseURL:  "https://api.exemplo.com",
		Token: &Token{
			AccessToken: "abc",
			ExpiresAt:   time.Now().Add(time.Hour).Round(time.Second),
			UserID:      "u-1",
		},
	}
	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Email != want.Email || got.Password != want.Password || got.BaseURL != want.BaseURL {
		t.Fatalf("campos divergentes: %+v", got)
	}
	if got.Token == nil || got.Token.AccessToken != "abc" || got.Token.UserID != "u-1" {
		t.Fatalf("token divergente: %+v", got.Token)
	}
	if !got.Token.ExpiresAt.Equal(want.Token.ExpiresAt) {
		t.Fatalf("expiresAt = %v, quer %v", got.Token.ExpiresAt, want.Token.ExpiresAt)
	}
}

func TestFileStore_Permissions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".breno")
	path := filepath.Join(dir, "config.json")
	s := NewFileStore(path)
	if err := s.Save(File{Email: "x"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat arquivo: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Fatalf("perm do arquivo = %o, quer 600", perm)
	}
	di, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if perm := di.Mode().Perm(); perm != 0o700 {
		t.Fatalf("perm do dir = %o, quer 700", perm)
	}
}

func TestFileStore_Clear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	s := NewFileStore(path)
	if err := s.Save(File{Email: "x"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("arquivo ainda existe: %v", err)
	}
	// Clear em arquivo inexistente não é erro.
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear idempotente: %v", err)
	}
}
