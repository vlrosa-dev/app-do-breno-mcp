package mcpserver

import (
	"context"
	"errors"
	"strings"
	"testing"

	"app-do-breno-agent/internal/breno"
	"app-do-breno-agent/internal/structs"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// fakeFinancas registra as chamadas recebidas e devolve respostas fixas.
type fakeFinancas struct {
	breno.Financas // embutido: métodos não sobrescritos entram em panic se usados
	novaMov        breno.NovaMovimentacao
	filtro         breno.FiltroMovimentacoes
	criarID        string
	criarErr       error
	listarResp     []breno.Movimentacao
}

func (f *fakeFinancas) CriarMovimentacao(_ context.Context, in breno.NovaMovimentacao) (string, error) {
	f.novaMov = in
	return f.criarID, f.criarErr
}

func (f *fakeFinancas) ListarMovimentacoes(_ context.Context, filtro breno.FiltroMovimentacoes) ([]breno.Movimentacao, error) {
	f.filtro = filtro
	return f.listarResp, nil
}

func textOf(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if res == nil || len(res.Content) == 0 {
		t.Fatal("resultado sem conteúdo")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("conteúdo não textual: %T", res.Content[0])
	}
	return tc.Text
}

func TestCreateTransactionTool_MapeiaEChamaServico(t *testing.T) {
	fake := &fakeFinancas{criarID: "mov-9"}
	h := HandleCreateTransactionTool(fake)

	res, _, err := h(context.Background(), nil, &structs.TransactionRequest{
		Tipo: "Diario", Valor: 12.5, Descricao: "Padaria", Data: "2026-09-03",
		TagIds: []string{"t1"},
	})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	if fake.novaMov.Tipo != breno.Diario || fake.novaMov.Valor != 12.5 || fake.novaMov.Descricao != "Padaria" {
		t.Fatalf("mapeamento errado: %+v", fake.novaMov)
	}
	if fake.novaMov.Recorrencia != nil {
		t.Fatal("não deveria ter recorrência")
	}
	if got := textOf(t, res); !strings.Contains(got, "mov-9") {
		t.Fatalf("texto = %q", got)
	}
}

func TestCreateTransactionTool_ErroServicoViraTexto(t *testing.T) {
	fake := &fakeFinancas{criarErr: errors.New("valor deve ser maior que zero")}
	h := HandleCreateTransactionTool(fake)
	res, _, err := h(context.Background(), nil, &structs.TransactionRequest{Tipo: "Diario"})
	if err != nil {
		t.Fatalf("handler não deveria propagar erro: %v", err)
	}
	if got := textOf(t, res); !strings.Contains(got, "maior que zero") {
		t.Fatalf("texto = %q", got)
	}
}

func TestCreateTransactionInSerieTool_PassaRecorrencia(t *testing.T) {
	fake := &fakeFinancas{}
	h := HandleCreateTransactionInSerieTool(fake)
	_, _, err := h(context.Background(), nil, &structs.TransactionInSerieRequest{
		Tipo: "Saida", Valor: 100, Descricao: "Aluguel", DataInicial: "2026-09-01",
		Frequencia: "mensalmente", NumeroRepeticoes: 12,
	})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	if fake.novaMov.Recorrencia == nil || fake.novaMov.Recorrencia.Repeticoes != 12 {
		t.Fatalf("recorrência = %+v", fake.novaMov.Recorrencia)
	}
	if fake.novaMov.Data != "2026-09-01" {
		t.Fatalf("data = %q", fake.novaMov.Data)
	}
}

func TestFetchTransactionTool_MapeiaFiltroEEstrutura(t *testing.T) {
	fake := &fakeFinancas{listarResp: []breno.Movimentacao{{MovimentacaoID: "m1", Descricao: "Padaria"}}}
	h := HandleFetchTransactionTool(fake)
	res, structured, err := h(context.Background(), nil, &structs.TransactionListRequest{
		Tipo: "diario", DataInicial: "2026-09-01", DataFinal: "2026-09-30",
	})
	if err != nil {
		t.Fatalf("handler: %v", err)
	}
	if fake.filtro.De != "2026-09-01" || fake.filtro.Ate != "2026-09-30" || fake.filtro.Tipo != "diario" {
		t.Fatalf("filtro = %+v", fake.filtro)
	}
	if got := textOf(t, res); !strings.Contains(got, "Padaria") {
		t.Fatalf("texto = %q", got)
	}
	if _, ok := structured.(listResult[breno.Movimentacao]); !ok {
		t.Fatalf("structuredContent = %T", structured)
	}
}
