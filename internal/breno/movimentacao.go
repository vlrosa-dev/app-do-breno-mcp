package breno

import (
	"context"
	"net/http"
	"net/url"

	"app-do-breno-agent/internal/structs"
)

// CriarMovimentacao cria uma movimentação única ou, quando in.Recorrencia != nil,
// uma série recorrente. Toda a validação de regra de negócio acontece aqui.
func (s *Service) CriarMovimentacao(ctx context.Context, in NovaMovimentacao) (string, error) {
	if err := in.normalize(); err != nil {
		return "", err
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return "", err
	}

	if in.Recorrencia == nil {
		payload := structs.TransactionRequest{
			Data:      in.Data,
			Descricao: in.Descricao,
			TagIds:    in.TagIDs,
			Tipo:      string(in.Tipo),
			User:      uid,
			Valor:     in.Valor,
		}
		var id string
		if err := s.client.RPC(ctx, "movimentacao_create_with_tags_v2", payload, &id); err != nil {
			return "", err
		}
		return id, nil
	}

	rec := in.Recorrencia
	payload := structs.TransactionInSerieRequest{
		DataInicial:              in.Data,
		Descricao:                in.Descricao,
		Frequencia:               string(rec.Frequencia),
		Infinita:                 rec.Infinita,
		NumeroRepeticoes:         rec.Repeticoes,
		PeriodicidadeRecorrencia: string(rec.Frequencia),
		QuantidadeRepeticoes:     rec.Repeticoes,
		TagIds:                   in.TagIDs,
		Tipo:                     string(in.Tipo),
		UserId:                   uid,
		Valor:                    in.Valor,
	}
	if err := s.client.RPC(ctx, "criar_movimentacoes_repetidas_v2", payload, nil); err != nil {
		return "", err
	}
	return "", nil
}

// AtualizarMovimentacao atualiza uma movimentação única pelo id.
func (s *Service) AtualizarMovimentacao(ctx context.Context, in AtualizacaoMovimentacao) error {
	if err := in.normalize(); err != nil {
		return err
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return err
	}
	payload := structs.TransactionUpdateRequest{
		Data:      in.Data,
		Descricao: in.Descricao,
		MovId:     in.ID,
		TagIds:    in.TagIDs,
		Tipo:      string(in.Tipo),
		User:      uid,
		Valor:     in.Valor,
	}
	return s.client.RPC(ctx, "movimentacao_update_with_tags_v2", payload, nil)
}

// DeletarMovimentacao remove uma movimentação única pelo id.
func (s *Service) DeletarMovimentacao(ctx context.Context, id string) error {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return err
	}
	q := url.Values{
		"id":      {eq(id)},
		"user_id": {eq(uid)},
	}
	return s.client.REST(ctx, http.MethodDelete, "/rest/v1/movimentacao", q, nil, nil)
}

// AtualizarSerie reancorá e atualiza os dados de uma série a partir de uma data.
func (s *Service) AtualizarSerie(ctx context.Context, in AtualizacaoSerie) error {
	if err := in.normalize(); err != nil {
		return err
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return err
	}
	payload := structs.TransactionInSerieUpdateRequest{
		DateMode:    "reanchor",
		Descricao:   in.Descricao,
		From:        in.De,
		NewDate:     in.NovaData,
		Recorrencia: in.RecorrenciaID,
		TagIds:      in.TagIDs,
		Tipo:        string(in.Tipo),
		User:        uid,
		Valor:       in.Valor,
	}
	return s.client.RPC(ctx, "serie_update_from_with_tags_v2", payload, nil)
}

// DeletarSerieDaData remove as ocorrências de uma série a partir de `de`.
func (s *Service) DeletarSerieDaData(ctx context.Context, recorrenciaID, de string) error {
	data, err := mustDate("data de origem (de)", de)
	if err != nil {
		return err
	}
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return err
	}
	payload := structs.TransactionInSerieDeleteRequest{
		From:          data,
		RecorrenciaID: recorrenciaID,
		UserId:        uid,
	}
	return s.client.RPC(ctx, "serie_delete_from", payload, nil)
}

// DeletarSerieToda remove todas as ocorrências de uma série.
func (s *Service) DeletarSerieToda(ctx context.Context, recorrenciaID string) error {
	uid, err := s.id.EnsureUserID(ctx)
	if err != nil {
		return err
	}
	payload := structs.TransactionInSerieDeleteAllRequest{
		RecorrenciaID: recorrenciaID,
		UserId:        uid,
	}
	return s.client.RPC(ctx, "serie_delete_all", payload, nil)
}
