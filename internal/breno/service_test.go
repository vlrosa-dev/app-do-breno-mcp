package breno

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

type restCall struct {
	method string
	path   string
	query  url.Values
	body   any
}

type fakeDoer struct {
	rest   []restCall
	rpc    map[string]any // fn -> body recebido
	rpcOut map[string]any // fn -> valor a escrever em out
	err    error
}

func newFakeDoer() *fakeDoer {
	return &fakeDoer{rpc: map[string]any{}, rpcOut: map[string]any{}}
}

func (f *fakeDoer) REST(_ context.Context, method, path string, query url.Values, body, out any) error {
	f.rest = append(f.rest, restCall{method, path, query, body})
	return f.err
}

func (f *fakeDoer) RPC(_ context.Context, fn string, body, out any) error {
	f.rpc[fn] = body
	if f.err != nil {
		return f.err
	}
	if v, ok := f.rpcOut[fn]; ok && out != nil {
		raw, _ := json.Marshal(v)
		return json.Unmarshal(raw, out)
	}
	return nil
}

type fakeID struct {
	uid string
	err error
}

func (f fakeID) EnsureUserID(context.Context) (string, error) { return f.uid, f.err }

func bodyField(t *testing.T, body any, jsonKey string) any {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	return m[jsonKey]
}

func TestCriarMovimentacaoUnica(t *testing.T) {
	doer := newFakeDoer()
	doer.rpcOut["movimentacao_create_with_tags_v2"] = "mov-1"
	svc := NewService(doer, fakeID{uid: "u-9"})

	id, err := svc.CriarMovimentacao(context.Background(), NovaMovimentacao{
		Tipo: "diario", Valor: 1.0, Descricao: "Padaria",
	})
	if err != nil {
		t.Fatalf("CriarMovimentacao: %v", err)
	}
	if id != "mov-1" {
		t.Fatalf("id = %q", id)
	}
	body := doer.rpc["movimentacao_create_with_tags_v2"]
	if got := bodyField(t, body, "p_user"); got != "u-9" {
		t.Fatalf("p_user = %v", got)
	}
	if got := bodyField(t, body, "p_tipo"); got != "Diario" {
		t.Fatalf("p_tipo = %v", got)
	}
	if got := bodyField(t, body, "p_data"); got == "" || got == nil {
		t.Fatalf("p_data vazio")
	}
	if _, criouSerie := doer.rpc["criar_movimentacoes_repetidas_v2"]; criouSerie {
		t.Fatal("não deveria ter criado série")
	}
}

func TestCriarMovimentacaoSerieInfinita(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-1"})

	_, err := svc.CriarMovimentacao(context.Background(), NovaMovimentacao{
		Tipo: "saida", Valor: 100, Descricao: "Aluguel",
		Recorrencia: &Recorrencia{Frequencia: "mensalmente", Infinita: true},
	})
	if err != nil {
		t.Fatalf("CriarMovimentacao série: %v", err)
	}
	body, ok := doer.rpc["criar_movimentacoes_repetidas_v2"]
	if !ok {
		t.Fatal("série não criada")
	}
	if got := bodyField(t, body, "p_numero_repeticoes"); got != float64(24) {
		t.Fatalf("p_numero_repeticoes = %v", got)
	}
	if got := bodyField(t, body, "p_frequencia"); got != "mensalmente" {
		t.Fatalf("p_frequencia = %v", got)
	}
}

func TestCriarMovimentacaoValidacaoNaoChamaCliente(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-1"})
	if _, err := svc.CriarMovimentacao(context.Background(), NovaMovimentacao{Tipo: "diario", Valor: 0, Descricao: "x"}); err == nil {
		t.Fatal("esperava erro de validação")
	}
	if len(doer.rpc) != 0 || len(doer.rest) != 0 {
		t.Fatal("cliente não deveria ter sido chamado")
	}
}

func TestDeletarMovimentacaoMontaQuery(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-2"})
	if err := svc.DeletarMovimentacao(context.Background(), "m-7"); err != nil {
		t.Fatalf("DeletarMovimentacao: %v", err)
	}
	if len(doer.rest) != 1 {
		t.Fatalf("esperava 1 chamada REST, veio %d", len(doer.rest))
	}
	c := doer.rest[0]
	if c.method != http.MethodDelete || c.path != "/rest/v1/movimentacao" {
		t.Fatalf("método/path = %s %s", c.method, c.path)
	}
	if c.query.Get("id") != "eq.m-7" || c.query.Get("user_id") != "eq.u-2" {
		t.Fatalf("query = %v", c.query)
	}
}

func TestListarMovimentacoesQuery(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-3"})
	_, err := svc.ListarMovimentacoes(context.Background(), FiltroMovimentacoes{
		Tipo: "diario", De: "2026-09-01", Ate: "2026-09-30",
	})
	if err != nil {
		t.Fatalf("ListarMovimentacoes: %v", err)
	}
	c := doer.rest[0]
	if c.method != http.MethodGet || c.path != "/rest/v1/view_movimentacao_all" {
		t.Fatalf("método/path = %s %s", c.method, c.path)
	}
	if got := c.query["data"]; len(got) != 2 || got[0] != "gte.2026-09-01" || got[1] != "lt.2026-09-30" {
		t.Fatalf("data query = %v", got)
	}
	if c.query.Get("tipo") != "eq.Diario" {
		t.Fatalf("tipo query = %q", c.query.Get("tipo"))
	}
}

func TestListarSaldosDefaultMeses(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-4"})
	if _, err := svc.ListarSaldos(context.Background(), "2026-09-01", 0); err != nil {
		t.Fatalf("ListarSaldos: %v", err)
	}
	if got := bodyField(t, doer.rpc["get_resumo_diario_web_optimized"], "p_quantidade_meses"); got != float64(1) {
		t.Fatalf("p_quantidade_meses = %v", got)
	}
}

func TestConsultasChamamRotaCerta(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-x"})
	ctx := context.Background()

	if _, err := svc.ObterMovimentacao(ctx, "m-1"); err != nil {
		t.Fatalf("ObterMovimentacao: %v", err)
	}
	if _, err := svc.ListarTags(ctx, "2026-09-01", "merc"); err != nil {
		t.Fatalf("ListarTags: %v", err)
	}
	if _, err := svc.ListarCartoes(ctx); err != nil {
		t.Fatalf("ListarCartoes: %v", err)
	}
	if _, err := svc.TotaisMes(ctx, "2026-09-01"); err != nil {
		t.Fatalf("TotaisMes: %v", err)
	}
	if _, err := svc.PerformanceAno(ctx, "2026"); err != nil {
		t.Fatalf("PerformanceAno: %v", err)
	}
	if _, err := svc.CustoDeVidaAno(ctx, "2026"); err != nil {
		t.Fatalf("CustoDeVidaAno: %v", err)
	}
	if _, err := svc.EconomiaAno(ctx, "2026"); err != nil {
		t.Fatalf("EconomiaAno: %v", err)
	}
	if _, err := svc.DiarioMedioAno(ctx, "2026"); err != nil {
		t.Fatalf("DiarioMedioAno: %v", err)
	}

	for _, fn := range []string{
		"get_totais_tag_mes_v3", "get_totais_mes", "get_performance_ano",
		"get_custo_vida_ano", "get_economizado_ano", "get_diario_medio_ano",
	} {
		if _, ok := doer.rpc[fn]; !ok {
			t.Errorf("RPC %q não chamada", fn)
		}
	}
	// ObterMovimentacao e ListarCartoes são GET REST.
	if len(doer.rest) != 2 {
		t.Fatalf("esperava 2 chamadas REST, veio %d", len(doer.rest))
	}
	if bodyField(t, doer.rpc["get_performance_ano"], "p_ano") != "2026" {
		t.Fatalf("p_ano não propagado")
	}
	if bodyField(t, doer.rpc["get_totais_tag_mes_v3"], "p_user_id") != "u-x" {
		t.Fatalf("p_user_id não propagado")
	}
}

func TestSeriesChamamRotaCerta(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-s"})
	ctx := context.Background()

	if err := svc.AtualizarSerie(ctx, AtualizacaoSerie{
		RecorrenciaID: "r1", De: "2026-09-01", NovaData: "2026-10-01",
		Tipo: "saida", Valor: 10, Descricao: "x",
	}); err != nil {
		t.Fatalf("AtualizarSerie: %v", err)
	}
	if err := svc.DeletarSerieDaData(ctx, "r1", "2026-10-01"); err != nil {
		t.Fatalf("DeletarSerieDaData: %v", err)
	}
	if err := svc.DeletarSerieToda(ctx, "r1"); err != nil {
		t.Fatalf("DeletarSerieToda: %v", err)
	}
	for _, fn := range []string{"serie_update_from_with_tags_v2", "serie_delete_from", "serie_delete_all"} {
		if _, ok := doer.rpc[fn]; !ok {
			t.Errorf("RPC %q não chamada", fn)
		}
	}
	if got := bodyField(t, doer.rpc["serie_update_from_with_tags_v2"], "p_date_mode"); got != "reanchor" {
		t.Fatalf("p_date_mode = %v", got)
	}
}

func TestAtualizarMovimentacaoValida(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{uid: "u-a"})
	err := svc.AtualizarMovimentacao(context.Background(), AtualizacaoMovimentacao{
		ID: "m1", Tipo: "diario", Valor: 5, Descricao: "y", Data: "2026-09-03",
	})
	if err != nil {
		t.Fatalf("AtualizarMovimentacao: %v", err)
	}
	if got := bodyField(t, doer.rpc["movimentacao_update_with_tags_v2"], "p_mov"); got != "m1" {
		t.Fatalf("p_mov = %v", got)
	}
	// sem data => erro de validação
	if err := svc.AtualizarMovimentacao(context.Background(), AtualizacaoMovimentacao{
		ID: "m1", Tipo: "diario", Valor: 5, Descricao: "y",
	}); err == nil {
		t.Fatal("esperava erro por data ausente")
	}
}

func TestIdentityErrorPropaga(t *testing.T) {
	doer := newFakeDoer()
	svc := NewService(doer, fakeID{err: errors.New("sem credenciais")})
	if _, err := svc.CriarMovimentacao(context.Background(), NovaMovimentacao{Tipo: "diario", Valor: 1, Descricao: "x"}); err == nil {
		t.Fatal("esperava erro de identidade")
	}
	if len(doer.rpc) != 0 {
		t.Fatal("não deveria chamar RPC sem uid")
	}
}
