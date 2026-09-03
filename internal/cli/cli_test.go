package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

// fakeAPI simula os endpoints Supabase usados pela CLI.
func fakeAPI(t *testing.T) (*httptest.Server, map[string]int) {
	t.Helper()
	hits := map[string]int{}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/v1/token", func(w http.ResponseWriter, r *http.Request) {
		hits["token"]++
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "acc-tok",
			"expires_at":   4102444800, // 2100-01-01
		})
	})
	mux.HandleFunc("/auth/v1/user", func(w http.ResponseWriter, r *http.Request) {
		hits["user"]++
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "user-1"})
	})
	mux.HandleFunc("/rest/v1/rpc/movimentacao_create_with_tags_v2", func(w http.ResponseWriter, r *http.Request) {
		hits["criar"]++
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"p_user":"user-1"`) {
			t.Errorf("payload sem p_user resolvido: %s", body)
		}
		_, _ = io.WriteString(w, `"mov-123"`)
	})
	mux.HandleFunc("/rest/v1/view_movimentacao_all", func(w http.ResponseWriter, r *http.Request) {
		hits["listar"]++
		_, _ = io.WriteString(w, `[{"movimentacao_id":"m1","data":"2026-09-03","valor":1.5,"descricao":"Padaria","tipo":"Diario"}]`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, hits
}

func run(t *testing.T, srv *httptest.Server, cfgPath string, args ...string) (string, string, error) {
	t.Helper()
	root := NewRootCmd()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	full := append([]string{"--config", cfgPath, "--base-url", srv.URL}, args...)
	root.SetArgs(full)
	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

func TestCLI_Help(t *testing.T) {
	srv, _ := fakeAPI(t)
	stdout, _, err := run(t, srv, filepath.Join(t.TempDir(), "c.json"), "--help")
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	for _, want := range []string{"diario", "saldos", "auth", "mov"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help não menciona %q", want)
		}
	}
}

func TestCLI_CriarDiario(t *testing.T) {
	srv, hits := fakeAPI(t)
	t.Setenv("API_USERNAME", "a@b.com")
	t.Setenv("API_PASSWORD", "pw")
	cfg := filepath.Join(t.TempDir(), "c.json")

	stdout, stderr, err := run(t, srv, cfg,
		"diario", "--valor", "1.00", "--descricao", "Padaria", "--recorrencia", "0")
	if err != nil {
		t.Fatalf("diario: %v (stderr=%s)", err, stderr)
	}
	if !strings.Contains(stdout, "mov-123") {
		t.Fatalf("stdout = %q", stdout)
	}
	if hits["criar"] != 1 || hits["token"] != 1 {
		t.Fatalf("hits = %v", hits)
	}
}

func TestCLI_CriarDiario_ValidacaoFaltandoFlag(t *testing.T) {
	srv, hits := fakeAPI(t)
	cfg := filepath.Join(t.TempDir(), "c.json")
	_, _, err := run(t, srv, cfg, "diario", "--descricao", "Padaria")
	if err == nil {
		t.Fatal("esperava erro por --valor faltando")
	}
	if hits["criar"] != 0 {
		t.Fatal("não deveria ter chamado a API")
	}
}

func TestCLI_ListarJSON(t *testing.T) {
	srv, _ := fakeAPI(t)
	t.Setenv("API_USERNAME", "a@b.com")
	t.Setenv("API_PASSWORD", "pw")
	cfg := filepath.Join(t.TempDir(), "c.json")

	stdout, stderr, err := run(t, srv, cfg,
		"--output", "json", "mov", "listar", "--tipo", "diario", "--de", "2026-09-01", "--ate", "2026-09-30")
	if err != nil {
		t.Fatalf("listar: %v (stderr=%s)", err, stderr)
	}
	var got []map[string]any
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("saída não é JSON: %v\n%s", err, stdout)
	}
	if len(got) != 1 || got[0]["descricao"] != "Padaria" {
		t.Fatalf("payload = %v", got)
	}
}

func TestCLI_ListarTabela(t *testing.T) {
	srv, _ := fakeAPI(t)
	t.Setenv("API_USERNAME", "a@b.com")
	t.Setenv("API_PASSWORD", "pw")
	cfg := filepath.Join(t.TempDir(), "c.json")

	stdout, _, err := run(t, srv, cfg,
		"mov", "listar", "--tipo", "diario", "--de", "2026-09-01", "--ate", "2026-09-30")
	if err != nil {
		t.Fatalf("listar: %v", err)
	}
	if !strings.Contains(stdout, "DESCRIÇÃO") || !strings.Contains(stdout, "Padaria") {
		t.Fatalf("tabela inesperada:\n%s", stdout)
	}
}

func TestCLI_AuthStatusSemLogin(t *testing.T) {
	srv, _ := fakeAPI(t)
	cfg := filepath.Join(t.TempDir(), "c.json")
	stdout, _, err := run(t, srv, cfg, "auth", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(stdout, "não autenticado") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCLI_AuthLoginPersiste(t *testing.T) {
	srv, hits := fakeAPI(t)
	cfg := filepath.Join(t.TempDir(), "breno", "config.json")

	stdout, stderr, err := run(t, srv, cfg, "auth", "login", "--email", "a@b.com", "--senha", "pw")
	if err != nil {
		t.Fatalf("login: %v (stderr=%s)", err, stderr)
	}
	if !strings.Contains(stdout, "autenticado como a@b.com") {
		t.Fatalf("stdout = %q", stdout)
	}
	if hits["token"] != 1 {
		t.Fatalf("hits token = %d", hits["token"])
	}
	// Segundo comando reaproveita o token do arquivo (sem novo login).
	if _, _, err := run(t, srv, cfg, "mov", "listar", "--tipo", "diario", "--de", "2026-09-01", "--ate", "2026-09-02"); err != nil {
		t.Fatalf("listar após login: %v", err)
	}
	if hits["token"] != 1 {
		t.Fatalf("relogin indevido: hits token = %d", hits["token"])
	}
}
