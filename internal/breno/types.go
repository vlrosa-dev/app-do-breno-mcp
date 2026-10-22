// Package breno concentra as regras de negócio do App do Breno num único lugar,
// independente do transporte. O servidor MCP depende apenas da interface
// Financas e dos tipos deste pacote — nunca dos structs de wire.
package breno

import "app-do-breno-agent/internal/structs"

// TipoMovimentacao é o tipo de uma movimentação financeira.
type TipoMovimentacao string

const (
	Diario   TipoMovimentacao = "Diario"
	Saida    TipoMovimentacao = "Saida"
	Cartao   TipoMovimentacao = "Cartao"
	Entrada  TipoMovimentacao = "Entrada"
	Economia TipoMovimentacao = "Economia"
)

// Frequencia é a periodicidade de uma série recorrente.
type Frequencia string

const (
	Mensalmente  Frequencia = "mensalmente"
	Semanalmente Frequencia = "semanalmente"
	Diariamente  Frequencia = "diariamente"
	Parcelado    Frequencia = "parcelado"
)

// Recorrencia descreve a repetição de uma movimentação. Presença de um ponteiro
// não-nil em NovaMovimentacao transforma a criação numa série.
type Recorrencia struct {
	Frequencia Frequencia
	Repeticoes int  // > 0; ignorado quando Infinita
	Infinita   bool // true => série "sem fim" (limitada internamente a 24)
}

// NovaMovimentacao é a entrada de CriarMovimentacao (única ou em série).
type NovaMovimentacao struct {
	Tipo        TipoMovimentacao
	Valor       float64
	Descricao   string
	Data        string // YYYY-MM-DD; vazio => hoje
	TagIDs      []string
	Recorrencia *Recorrencia // nil => movimentação única
}

// AtualizacaoMovimentacao é a entrada de AtualizarMovimentacao.
type AtualizacaoMovimentacao struct {
	ID        string
	Tipo      TipoMovimentacao
	Valor     float64
	Descricao string
	Data      string // YYYY-MM-DD
	TagIDs    []string
}

// AtualizacaoSerie é a entrada de AtualizarSerie.
type AtualizacaoSerie struct {
	RecorrenciaID string
	De            string // YYYY-MM-DD — âncora atual
	NovaData      string // YYYY-MM-DD — nova âncora
	Tipo          TipoMovimentacao
	Valor         float64
	Descricao     string
	TagIDs        []string
}

// FiltroMovimentacoes é a entrada de ListarMovimentacoes.
type FiltroMovimentacoes struct {
	Tipo   TipoMovimentacao
	De     string // YYYY-MM-DD
	Ate    string // YYYY-MM-DD
	Select string // colunas PostgREST; vazio => "*"
}

// Tipos de resposta (read models). Por ora são aliases dos structs de wire —
// o servidor MCP importa só o pacote breno, e a camada pode endurecer
// depois sem quebrar assinaturas.
type (
	Movimentacao   = structs.TransactionListResponse
	Saldo          = structs.ResumoSaldosResponse
	Tag            = structs.TagListResponse
	CartaoCredito  = structs.ListCartaoResponse
	TotaisMes      = structs.TotaisMesListResponse
	PerformanceAno = structs.TotaisAnoPerformaceResponse
	CustoDeVidaAno = structs.TotaisAnoCustoDeVidaResponse
	EconomiaAno    = structs.TotaisAnoEconomizadoResponse
	DiarioMedioAno = structs.TotaisDiarioMedioAnoResponse
)
