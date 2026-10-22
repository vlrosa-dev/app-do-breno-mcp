//go:build integration

// Package integration exercita as rotas reais da API (Supabase) passando pela
// mesma pilha usada pelo breno-mcp: auth.Manager -> api.Client -> breno.Service
// -> servidor MCP.
//
// Rodar com `make test-integration`. Requer API_USERNAME/API_PASSWORD (ambiente
// ou .env na raiz). As rotas de escrita só rodam com BRENO_IT_WRITE=1 — elas
// criam movimentações de R$ 0,01 marcadas com "[it-mcp]" e as removem ao final.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/auth"
	"app-do-breno-agent/internal/breno"
	"app-do-breno-agent/internal/config"
	"app-do-breno-agent/internal/mcpserver"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const marcador = "[it-mcp]"

// newService monta a pilha real com um Store temporário, para não tocar em
// ~/.breno/config.json.
func newService(t *testing.T) *breno.Service {
	t.Helper()
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	cfg := config.Load()
	if cfg.Email == "" || cfg.Password == "" {
		t.Skip("API_USERNAME/API_PASSWORD não definidos — pulando testes de integração")
	}

	manager, err := auth.NewManager(auth.Options{
		BaseURL:     cfg.APIBaseURL,
		APIKey:      cfg.APIKey,
		Store:       auth.NewFileStore(filepath.Join(t.TempDir(), "config.json")),
		Credentials: auth.Credentials{Email: cfg.Email, Password: cfg.Password},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	client := api.NewClient(api.ClientOptions{
		BaseURL: cfg.APIBaseURL,
		APIKey:  cfg.APIKey,
		Auth:    manager,
	})
	return breno.NewService(client, manager)
}

func ctxT(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func inicioDoMes(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

// Consultas (somente leitura)

func TestRotasDeConsulta(t *testing.T) {
	svc := newService(t)
	ctx := ctxT(t)

	hoje := time.Now()
	mes := inicioDoMes(hoje).Format("2006-01-02")
	ate := inicioDoMes(hoje).AddDate(0, 1, 0).Format("2006-01-02")
	ano := hoje.Format("2006")

	for _, tipo := range []breno.TipoMovimentacao{breno.Diario, breno.Saida, breno.Cartao, breno.Entrada, breno.Economia} {
		t.Run("view_movimentacao_all/"+string(tipo), func(t *testing.T) {
			movs, err := svc.ListarMovimentacoes(ctx, breno.FiltroMovimentacoes{Tipo: tipo, De: mes, Ate: ate})
			if err != nil {
				t.Fatalf("ListarMovimentacoes: %v", err)
			}
			for _, m := range movs {
				if m.Tipo != string(tipo) {
					t.Fatalf("tipo = %q, quer %q", m.Tipo, tipo)
				}
				if m.Data < mes || m.Data >= ate {
					t.Fatalf("data %s fora de [%s, %s)", m.Data, mes, ate)
				}
			}
			t.Logf("%d movimentações", len(movs))

			if len(movs) == 0 {
				return
			}
			got, err := svc.ObterMovimentacao(ctx, movs[0].MovimentacaoID)
			if err != nil {
				t.Fatalf("ObterMovimentacao: %v", err)
			}
			if len(got) == 0 || got[0].MovimentacaoID != movs[0].MovimentacaoID {
				t.Fatalf("ObterMovimentacao(%s) = %+v", movs[0].MovimentacaoID, got)
			}
		})
	}

	t.Run("rpc/get_resumo_diario_web_optimized", func(t *testing.T) {
		out, err := svc.ListarSaldos(ctx, mes, 3)
		if err != nil {
			t.Fatalf("ListarSaldos: %v", err)
		}
		t.Logf("%d saldos", len(out))
	})
	t.Run("rpc/get_totais_tag_mes_v3", func(t *testing.T) {
		out, err := svc.ListarTags(ctx, mes, "")
		if err != nil {
			t.Fatalf("ListarTags: %v", err)
		}
		t.Logf("%d tags", len(out))
	})
	t.Run("rest/cartao", func(t *testing.T) {
		out, err := svc.ListarCartoes(ctx)
		if err != nil {
			t.Fatalf("ListarCartoes: %v", err)
		}
		t.Logf("%d cartões", len(out))
	})
	t.Run("rpc/get_totais_mes", func(t *testing.T) {
		if _, err := svc.TotaisMes(ctx, mes); err != nil {
			t.Fatalf("TotaisMes: %v", err)
		}
	})
	t.Run("rpc/get_performance_ano", func(t *testing.T) {
		if _, err := svc.PerformanceAno(ctx, ano); err != nil {
			t.Fatalf("PerformanceAno: %v", err)
		}
	})
	t.Run("rpc/get_custo_vida_ano", func(t *testing.T) {
		if _, err := svc.CustoDeVidaAno(ctx, ano); err != nil {
			t.Fatalf("CustoDeVidaAno: %v", err)
		}
	})
	t.Run("rpc/get_economizado_ano", func(t *testing.T) {
		if _, err := svc.EconomiaAno(ctx, ano); err != nil {
			t.Fatalf("EconomiaAno: %v", err)
		}
	})
	t.Run("rpc/get_diario_medio_ano", func(t *testing.T) {
		if _, err := svc.DiarioMedioAno(ctx, ano); err != nil {
			t.Fatalf("DiarioMedioAno: %v", err)
		}
	})
}

// Escrita (BRENO_IT_WRITE=1)

func requireWrite(t *testing.T) {
	t.Helper()
	if os.Getenv("BRENO_IT_WRITE") != "1" {
		t.Skip("BRENO_IT_WRITE != 1 — pulando rotas de escrita")
	}
}

// buscarPorDescricao lista as movimentações do tipo no intervalo e filtra pela
// descrição exata.
func buscarPorDescricao(t *testing.T, ctx context.Context, svc *breno.Service, tipo breno.TipoMovimentacao, de, ate, desc string) []breno.Movimentacao {
	t.Helper()
	movs, err := svc.ListarMovimentacoes(ctx, breno.FiltroMovimentacoes{Tipo: tipo, De: de, Ate: ate})
	if err != nil {
		t.Fatalf("ListarMovimentacoes: %v", err)
	}
	var out []breno.Movimentacao
	for _, m := range movs {
		if m.Descricao == desc {
			out = append(out, m)
		}
	}
	return out
}

func TestRotasDeMovimentacaoUnica(t *testing.T) {
	requireWrite(t)
	svc := newService(t)
	ctx := ctxT(t)

	hoje := time.Now().Format("2006-01-02")
	desc := fmt.Sprintf("%s unica %d", marcador, time.Now().UnixNano())

	id, err := svc.CriarMovimentacao(ctx, breno.NovaMovimentacao{
		Tipo: breno.Diario, Valor: 0.01, Descricao: desc, Data: hoje,
	})
	if err != nil {
		t.Fatalf("CriarMovimentacao (movimentacao_create_with_tags_v2): %v", err)
	}
	if id == "" {
		t.Fatal("CriarMovimentacao devolveu id vazio")
	}
	t.Cleanup(func() {
		// Garante a limpeza mesmo se um passo intermediário falhar.
		_ = svc.DeletarMovimentacao(context.Background(), id)
	})

	got, err := svc.ObterMovimentacao(ctx, id)
	if err != nil || len(got) == 0 {
		t.Fatalf("ObterMovimentacao após criar: %v (%d itens)", err, len(got))
	}
	if got[0].Descricao != desc || got[0].Valor != 0.01 {
		t.Fatalf("criada com dados divergentes: %+v", got[0])
	}

	novaDesc := desc + " editada"
	err = svc.AtualizarMovimentacao(ctx, breno.AtualizacaoMovimentacao{
		ID: id, Tipo: breno.Diario, Valor: 0.02, Descricao: novaDesc, Data: hoje,
	})
	if err != nil {
		t.Fatalf("AtualizarMovimentacao (movimentacao_update_with_tags_v2): %v", err)
	}
	got, err = svc.ObterMovimentacao(ctx, id)
	if err != nil || len(got) == 0 {
		t.Fatalf("ObterMovimentacao após atualizar: %v (%d itens)", err, len(got))
	}
	if got[0].Descricao != novaDesc || got[0].Valor != 0.02 {
		t.Fatalf("atualização não refletida: %+v", got[0])
	}

	if err := svc.DeletarMovimentacao(ctx, id); err != nil {
		t.Fatalf("DeletarMovimentacao (DELETE /rest/v1/movimentacao): %v", err)
	}
	got, err = svc.ObterMovimentacao(ctx, id)
	if err != nil {
		t.Fatalf("ObterMovimentacao após deletar: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("movimentação ainda existe após deletar: %+v", got)
	}
}

func TestRotasDeSerie(t *testing.T) {
	requireWrite(t)
	svc := newService(t)
	ctx := ctxT(t)

	inicio := inicioDoMes(time.Now())
	de := inicio.Format("2006-01-02")
	ate := inicio.AddDate(0, 4, 0).Format("2006-01-02")
	desc := fmt.Sprintf("%s serie %d", marcador, time.Now().UnixNano())

	_, err := svc.CriarMovimentacao(ctx, breno.NovaMovimentacao{
		Tipo: breno.Saida, Valor: 0.01, Descricao: desc, Data: de,
		Recorrencia: &breno.Recorrencia{Frequencia: breno.Mensalmente, Repeticoes: 3},
	})
	if err != nil {
		t.Fatalf("CriarMovimentacao em série (criar_movimentacoes_repetidas_v2): %v", err)
	}

	criadas := buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, desc)
	if len(criadas) == 0 {
		t.Fatal("série criada mas nenhuma ocorrência encontrada")
	}
	recID, _ := criadas[0].RecorrenciaID.(string)
	if recID == "" {
		t.Fatalf("ocorrência sem recorrencia_id: %+v", criadas[0])
	}
	t.Cleanup(func() {
		_ = svc.DeletarSerieToda(context.Background(), recID)
	})
	t.Logf("série %s com %d ocorrências", recID, len(criadas))

	// Atualiza a partir da 2ª ocorrência, mantendo a mesma âncora.
	segunda := inicio.AddDate(0, 1, 0).Format("2006-01-02")
	descEditada := desc + " editada"
	err = svc.AtualizarSerie(ctx, breno.AtualizacaoSerie{
		RecorrenciaID: recID, De: segunda, NovaData: segunda,
		Tipo: breno.Saida, Valor: 0.02, Descricao: descEditada,
	})
	if err != nil {
		t.Fatalf("AtualizarSerie (serie_update_from_with_tags_v2): %v", err)
	}
	if editadas := buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, descEditada); len(editadas) == 0 {
		t.Fatal("nenhuma ocorrência refletiu a atualização da série")
	}

	// Remove da 3ª ocorrência em diante.
	terceira := inicio.AddDate(0, 2, 0).Format("2006-01-02")
	if err := svc.DeletarSerieDaData(ctx, recID, terceira); err != nil {
		t.Fatalf("DeletarSerieDaData (serie_delete_from): %v", err)
	}
	restantes := append(
		buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, desc),
		buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, descEditada)...,
	)
	for _, m := range restantes {
		if m.Data >= terceira {
			t.Fatalf("ocorrência de %s deveria ter sido removida", m.Data)
		}
	}

	if err := svc.DeletarSerieToda(ctx, recID); err != nil {
		t.Fatalf("DeletarSerieToda (serie_delete_all): %v", err)
	}
	restantes = append(
		buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, desc),
		buscarPorDescricao(t, ctx, svc, breno.Saida, de, ate, descEditada)...,
	)
	if len(restantes) != 0 {
		t.Fatalf("série ainda tem %d ocorrências após serie_delete_all", len(restantes))
	}
}

// Servidor MCP ponta a ponta (somente leitura)

func TestServidorMCP(t *testing.T) {
	svc := newService(t)
	ctx := ctxT(t)

	server, err := mcpserver.NewMcpServer(nil, svc)
	if err != nil {
		t.Fatalf("NewMcpServer: %v", err)
	}
	serverT, clientT := mcp.NewInMemoryTransports()
	ss, err := server.Connect(ctx, serverT, nil)
	if err != nil {
		t.Fatalf("server.Connect: %v", err)
	}
	t.Cleanup(func() { _ = ss.Close() })

	client := mcp.NewClient(&mcp.Implementation{Name: "it-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if len(tools.Tools) != 17 {
		t.Fatalf("esperava 17 tools, veio %d", len(tools.Tools))
	}

	// Sem campos de usuário: o servidor usa o id do usuário autenticado.
	mes := inicioDoMes(time.Now()).Format("2006-01-02")
	chamadas := map[string]map[string]any{
		"listar_cartoes_de_creditos": {},
		"listar_saldos":              {"p_data_mes": mes, "p_quantidade_meses": 1},
		"listar_totais_mes":          {"mes_input": mes},
	}
	for nome, args := range chamadas {
		t.Run(nome, func(t *testing.T) {
			res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: nome, Arguments: args})
			if err != nil {
				t.Fatalf("CallTool: %v", err)
			}
			if len(res.Content) == 0 {
				t.Fatalf("resultado sem conteúdo: %+v", res)
			}
			text := res.Content[0].(*mcp.TextContent).Text
			if res.IsError || strings.HasPrefix(text, "Unable to") || !json.Valid([]byte(text)) {
				t.Fatalf("resposta inesperada: %s", text)
			}
		})
	}
}
