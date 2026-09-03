// Package cli implementa a ferramenta de linha de comando `breno`, que
// reaproveita a mesma camada de regras (internal/breno) do servidor MCP.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/auth"
	"app-do-breno-agent/internal/breno"
	"app-do-breno-agent/internal/config"

	"github.com/spf13/cobra"
)

type globalFlags struct {
	output     string
	configPath string
	baseURL    string
	apiKey     string
}

// deps carrega dependências resolvidas de forma preguiçosa (só quando um
// comando realmente precisa falar com a API).
type deps struct {
	g       *globalFlags
	manager *auth.Manager
	svc     breno.Financas
	baseURL string
	apiKey  string
}

// NewRootCmd monta a árvore de comandos.
func NewRootCmd() *cobra.Command {
	g := &globalFlags{}
	d := &deps{g: g}

	root := &cobra.Command{
		Use:           "breno",
		Short:         "CLI do App do Breno — movimentações e consultas financeiras",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	pf := root.PersistentFlags()
	pf.StringVarP(&g.output, "output", "o", "table", "formato de saída: table|json")
	pf.StringVar(&g.configPath, "config", "", "caminho do arquivo de config (default ~/.breno/config.json)")
	pf.StringVar(&g.baseURL, "base-url", "", "override da URL base da API")
	pf.StringVar(&g.apiKey, "api-key", "", "override da API key")

	root.AddCommand(newMovimentacaoCmds(d)...)
	root.AddCommand(
		newMovGroupCmd(d),
		newSerieGroupCmd(d),
		newTotaisGroupCmd(d),
		newAuthGroupCmd(d),
	)
	root.AddCommand(newConsultaCmds(d)...)
	return root
}

// Execute é o ponto de entrada usado por cmd/breno-cli.
func Execute() int {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		return 1
	}
	return 0
}

// resolveConfigPath devolve o caminho do arquivo de config efetivo.
func (d *deps) resolveConfigPath() (string, error) {
	if d.g.configPath != "" {
		return d.g.configPath, nil
	}
	return auth.DefaultPath()
}

// authManager monta (uma vez) o auth.Manager com a precedência
// flag > env > arquivo > default.
func (d *deps) authManager() (*auth.Manager, error) {
	if d.manager != nil {
		return d.manager, nil
	}
	path, err := d.resolveConfigPath()
	if err != nil {
		return nil, err
	}
	store := auth.NewFileStore(path)
	f, err := store.Load()
	if err != nil {
		return nil, err
	}

	d.baseURL = firstNonEmpty(d.g.baseURL, os.Getenv("API_BASE_URL"), f.BaseURL, config.DefaultAPIBaseURL)
	d.apiKey = firstNonEmpty(d.g.apiKey, os.Getenv("API_KEY"), f.APIKey, config.DefaultPublicAPIKey)
	creds := auth.Credentials{
		Email:    firstNonEmpty(os.Getenv("API_USERNAME"), f.Email),
		Password: firstNonEmpty(os.Getenv("API_PASSWORD"), f.Password),
	}

	m, err := auth.NewManager(auth.Options{
		BaseURL:     d.baseURL,
		APIKey:      d.apiKey,
		Store:       store,
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	d.manager = m
	return m, nil
}

// service monta (uma vez) o breno.Service pronto para uso.
func (d *deps) service() (breno.Financas, error) {
	if d.svc != nil {
		return d.svc, nil
	}
	m, err := d.authManager()
	if err != nil {
		return nil, err
	}
	client := api.NewClient(api.ClientOptions{BaseURL: d.baseURL, APIKey: d.apiKey, Auth: m})
	d.svc = breno.NewService(client, m)
	return d.svc, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// cmdContext devolve o context do comando (ou Background).
func cmdContext(cmd *cobra.Command) context.Context {
	if ctx := cmd.Context(); ctx != nil {
		return ctx
	}
	return context.Background()
}

func out(cmd *cobra.Command) io.Writer { return cmd.OutOrStdout() }
