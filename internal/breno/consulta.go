package breno

import (
	"context"
	"net/http"
	"net/url"

	"github.com/vlrosa-dev/app-do-breno-mcp/internal/structs"
)

// ListarMovimentacoes filtra por tipo e período [de, ate).
func (s *Service) ListarMovimentacoes(ctx context.Context, f FiltroMovimentacoes) ([]Movimentacao, error) {
	tipo, de, ate, sel, err := f.resolve()
	if err != nil {
		return nil, err
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{
		"select":  {sel},
		"user_id": {eq(uid)},
		"tipo":    {eq(string(tipo))},
		"data":    {"gte." + de, "lt." + ate},
		"order":   {"data.desc,created_at.desc"},
	}
	var out []Movimentacao
	if err := s.client.REST(ctx, http.MethodGet, "/rest/v1/view_movimentacao_all", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ObterMovimentacao busca uma movimentação (e suas ocorrências) pelo id.
func (s *Service) ObterMovimentacao(ctx context.Context, id string) ([]Movimentacao, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{
		"select":          {"*"},
		"movimentacao_id": {eq(id)},
		"user_id":         {eq(uid)},
	}
	var out []Movimentacao
	if err := s.client.REST(ctx, http.MethodGet, "/rest/v1/view_movimentacao_all", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListarSaldos devolve os saldos por tipo a partir do mês informado.
func (s *Service) ListarSaldos(ctx context.Context, mes string, meses int) ([]Saldo, error) {
	if meses <= 0 {
		meses = 1
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return nil, err
	}
	payload := structs.ResumoSaldosRequest{
		DataMes:         mes,
		QuantidadeMeses: int64(meses),
		UserID:          uid,
	}
	var out []Saldo
	if err := s.client.RPC(ctx, "get_resumo_diario_web_optimized", payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListarTags devolve as tags e seus totais no mês (busca opcional).
func (s *Service) ListarTags(ctx context.Context, mes, busca string) ([]Tag, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return nil, err
	}
	payload := structs.TagListRequest{
		DataMes: mes,
		Search:  busca,
		UserId:  uid,
	}
	var out []Tag
	if err := s.client.RPC(ctx, "get_totais_tag_mes_v3", payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListarCartoes devolve os cartões de crédito cadastrados.
func (s *Service) ListarCartoes(ctx context.Context) ([]CartaoCredito, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{
		"select":  {"*"},
		"user_id": {eq(uid)},
		"order":   {"nome.asc"},
	}
	var out []CartaoCredito
	if err := s.client.REST(ctx, http.MethodGet, "/rest/v1/cartao", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TotaisMes devolve os indicadores de um mês (formato YYYY-MM-DD ou YYYY-MM).
func (s *Service) TotaisMes(ctx context.Context, mes string) (TotaisMes, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return TotaisMes{}, err
	}
	payload := structs.TotaisMesListRequest{MesInput: mes, UserInput: uid}
	var out TotaisMes
	if err := s.client.RPC(ctx, "get_totais_mes", payload, &out); err != nil {
		return TotaisMes{}, err
	}
	return out, nil
}

// PerformanceAno devolve a performance mês a mês no ano.
func (s *Service) PerformanceAno(ctx context.Context, ano string) (PerformanceAno, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return PerformanceAno{}, err
	}
	payload := structs.TotaisAnoPerformanceRequest{Ano: ano, UserId: uid}
	var out PerformanceAno
	if err := s.client.RPC(ctx, "get_performance_ano", payload, &out); err != nil {
		return PerformanceAno{}, err
	}
	return out, nil
}

// CustoDeVidaAno devolve o custo de vida mês a mês no ano.
func (s *Service) CustoDeVidaAno(ctx context.Context, ano string) (CustoDeVidaAno, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return CustoDeVidaAno{}, err
	}
	payload := structs.TotaisAnoCustoDeVidaRequest{Ano: ano, UserId: uid}
	var out CustoDeVidaAno
	if err := s.client.RPC(ctx, "get_custo_vida_ano", payload, &out); err != nil {
		return CustoDeVidaAno{}, err
	}
	return out, nil
}

// EconomiaAno devolve o total economizado mês a mês no ano.
func (s *Service) EconomiaAno(ctx context.Context, ano string) (EconomiaAno, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return EconomiaAno{}, err
	}
	payload := structs.TotaisAnoEconomizadoRequest{Ano: ano, UserId: uid}
	var out EconomiaAno
	if err := s.client.RPC(ctx, "get_economizado_ano", payload, &out); err != nil {
		return EconomiaAno{}, err
	}
	return out, nil
}

// DiarioMedioAno devolve o gasto diário médio mês a mês no ano.
func (s *Service) DiarioMedioAno(ctx context.Context, ano string) (DiarioMedioAno, error) {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return DiarioMedioAno{}, err
	}
	payload := structs.TotaisDiarioMedioAnoRequest{Ano: ano, UserId: uid}
	var out DiarioMedioAno
	if err := s.client.RPC(ctx, "get_diario_medio_ano", payload, &out); err != nil {
		return DiarioMedioAno{}, err
	}
	return out, nil
}
