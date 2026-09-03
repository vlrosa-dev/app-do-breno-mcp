package cli

import (
	"fmt"

	"app-do-breno-agent/internal/breno"

	"github.com/spf13/cobra"
)

type criarFlags struct {
	valor       float64
	descricao   string
	data        string
	tags        []string
	recorrencia int
	frequencia  string
	infinita    bool
}

func bindCriarFlags(cmd *cobra.Command, f *criarFlags) {
	fl := cmd.Flags()
	fl.Float64Var(&f.valor, "valor", 0, "valor da movimentação (obrigatório)")
	fl.StringVar(&f.descricao, "descricao", "", "descrição da movimentação (obrigatório)")
	fl.StringVar(&f.data, "data", "", "data YYYY-MM-DD (default: hoje)")
	fl.StringSliceVar(&f.tags, "tags", nil, "IDs de tags separados por vírgula")
	fl.IntVar(&f.recorrencia, "recorrencia", 0, "nº de repetições; 0 = movimentação única")
	fl.StringVar(&f.frequencia, "frequencia", "mensalmente", "mensalmente|semanalmente|diariamente|parcelado (quando --recorrencia > 0)")
	fl.BoolVar(&f.infinita, "infinita", false, "série sem fim (limitada internamente a 24)")
	_ = cmd.MarkFlagRequired("valor")
	_ = cmd.MarkFlagRequired("descricao")
}

func (f criarFlags) toNova(tipo breno.TipoMovimentacao) breno.NovaMovimentacao {
	nova := breno.NovaMovimentacao{
		Tipo:      tipo,
		Valor:     f.valor,
		Descricao: f.descricao,
		Data:      f.data,
		TagIDs:    f.tags,
	}
	if f.recorrencia > 0 || f.infinita {
		nova.Recorrencia = &breno.Recorrencia{
			Frequencia: breno.Frequencia(f.frequencia),
			Repeticoes: f.recorrencia,
			Infinita:   f.infinita,
		}
	}
	return nova
}

// newMovimentacaoCmds cria os comandos de atalho por tipo (diario, saida, ...).
func newMovimentacaoCmds(d *deps) []*cobra.Command {
	tipos := []struct {
		use  string
		tipo breno.TipoMovimentacao
		help string
	}{
		{"diario", breno.Diario, "Registra um gasto do dia a dia"},
		{"saida", breno.Saida, "Registra uma saída fixa mensal"},
		{"cartao", breno.Cartao, "Registra uma despesa em cartão de crédito"},
		{"entrada", breno.Entrada, "Registra uma receita recebida"},
		{"economia", breno.Economia, "Registra um valor economizado/poupado"},
	}

	cmds := make([]*cobra.Command, 0, len(tipos))
	for _, tt := range tipos {
		tt := tt
		f := &criarFlags{}
		cmd := &cobra.Command{
			Use:   tt.use,
			Short: tt.help,
			RunE: func(cmd *cobra.Command, _ []string) error {
				svc, err := d.service()
				if err != nil {
					return err
				}
				id, err := svc.CriarMovimentacao(cmdContext(cmd), f.toNova(tt.tipo))
				if err != nil {
					return err
				}
				msg := "movimentação criada"
				if id != "" {
					msg = fmt.Sprintf("movimentação criada: %s", id)
				}
				return okMessage(cmd, d.g, msg)
			},
		}
		bindCriarFlags(cmd, f)
		cmds = append(cmds, cmd)
	}
	return cmds
}

// newMovGroupCmd cria o grupo `mov` (listar/obter/atualizar/deletar).
func newMovGroupCmd(d *deps) *cobra.Command {
	group := &cobra.Command{Use: "mov", Short: "Consulta e edição de movimentações"}

	var (
		fTipo, fDe, fAte, fSelect string
	)
	listar := &cobra.Command{
		Use:   "listar",
		Short: "Lista movimentações por tipo e período",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			movs, err := svc.ListarMovimentacoes(cmdContext(cmd), breno.FiltroMovimentacoes{
				Tipo: breno.TipoMovimentacao(fTipo), De: fDe, Ate: fAte, Select: fSelect,
			})
			if err != nil {
				return err
			}
			return render(cmd, d.g, movs)
		},
	}
	listar.Flags().StringVar(&fTipo, "tipo", "", "tipo (obrigatório)")
	listar.Flags().StringVar(&fDe, "de", "", "data inicial YYYY-MM-DD (obrigatório)")
	listar.Flags().StringVar(&fAte, "ate", "", "data final YYYY-MM-DD (obrigatório)")
	listar.Flags().StringVar(&fSelect, "select", "", "colunas PostgREST (default *)")
	_ = listar.MarkFlagRequired("tipo")
	_ = listar.MarkFlagRequired("de")
	_ = listar.MarkFlagRequired("ate")

	var obterID string
	obter := &cobra.Command{
		Use:   "obter",
		Short: "Busca uma movimentação pelo ID",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			movs, err := svc.ObterMovimentacao(cmdContext(cmd), obterID)
			if err != nil {
				return err
			}
			return render(cmd, d.g, movs)
		},
	}
	obter.Flags().StringVar(&obterID, "id", "", "ID da movimentação (obrigatório)")
	_ = obter.MarkFlagRequired("id")

	var (
		upID, upTipo, upDescricao, upData string
		upValor                           float64
		upTags                            []string
	)
	atualizar := &cobra.Command{
		Use:   "atualizar",
		Short: "Atualiza uma movimentação única",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			err = svc.AtualizarMovimentacao(cmdContext(cmd), breno.AtualizacaoMovimentacao{
				ID: upID, Tipo: breno.TipoMovimentacao(upTipo), Valor: upValor,
				Descricao: upDescricao, Data: upData, TagIDs: upTags,
			})
			if err != nil {
				return err
			}
			return okMessage(cmd, d.g, "movimentação atualizada")
		},
	}
	af := atualizar.Flags()
	af.StringVar(&upID, "id", "", "ID da movimentação (obrigatório)")
	af.StringVar(&upTipo, "tipo", "", "tipo (obrigatório)")
	af.Float64Var(&upValor, "valor", 0, "novo valor (obrigatório)")
	af.StringVar(&upDescricao, "descricao", "", "nova descrição (obrigatório)")
	af.StringVar(&upData, "data", "", "nova data YYYY-MM-DD (obrigatório)")
	af.StringSliceVar(&upTags, "tags", nil, "IDs de tags")
	for _, r := range []string{"id", "tipo", "valor", "descricao", "data"} {
		_ = atualizar.MarkFlagRequired(r)
	}

	var delID string
	deletar := &cobra.Command{
		Use:   "deletar",
		Short: "Remove uma movimentação única",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			if err := svc.DeletarMovimentacao(cmdContext(cmd), delID); err != nil {
				return err
			}
			return okMessage(cmd, d.g, "movimentação removida")
		},
	}
	deletar.Flags().StringVar(&delID, "id", "", "ID da movimentação (obrigatório)")
	_ = deletar.MarkFlagRequired("id")

	group.AddCommand(listar, obter, atualizar, deletar)
	return group
}
