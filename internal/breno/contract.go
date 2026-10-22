package breno

import "context"

// Financas é o contrato de domínio consumido pelo servidor MCP.
// A implementação concreta é *Service.
type Financas interface {
	// Movimentações
	CriarMovimentacao(ctx context.Context, in NovaMovimentacao) (string, error)
	AtualizarMovimentacao(ctx context.Context, in AtualizacaoMovimentacao) error
	DeletarMovimentacao(ctx context.Context, id string) error

	// Séries recorrentes
	AtualizarSerie(ctx context.Context, in AtualizacaoSerie) error
	DeletarSerieDaData(ctx context.Context, recorrenciaID, de string) error
	DeletarSerieToda(ctx context.Context, recorrenciaID string) error

	// Consultas
	ListarMovimentacoes(ctx context.Context, f FiltroMovimentacoes) ([]Movimentacao, error)
	ObterMovimentacao(ctx context.Context, id string) ([]Movimentacao, error)
	ListarSaldos(ctx context.Context, mes string, meses int) ([]Saldo, error)
	ListarTags(ctx context.Context, mes, busca string) ([]Tag, error)
	ListarCartoes(ctx context.Context) ([]CartaoCredito, error)
	TotaisMes(ctx context.Context, mes string) (TotaisMes, error)
	PerformanceAno(ctx context.Context, ano string) (PerformanceAno, error)
	CustoDeVidaAno(ctx context.Context, ano string) (CustoDeVidaAno, error)
	EconomiaAno(ctx context.Context, ano string) (EconomiaAno, error)
	DiarioMedioAno(ctx context.Context, ano string) (DiarioMedioAno, error)
}
