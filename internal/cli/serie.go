package cli

import (
	"app-do-breno-agent/internal/breno"

	"github.com/spf13/cobra"
)

// newSerieGroupCmd cria o grupo `serie` (atualizar/deletar/deletar-tudo).
func newSerieGroupCmd(d *deps) *cobra.Command {
	group := &cobra.Command{Use: "serie", Short: "Edição de séries recorrentes"}

	var (
		upID, upDe, upNova, upTipo, upDescricao string
		upValor                                 float64
		upTags                                  []string
	)
	atualizar := &cobra.Command{
		Use:   "atualizar",
		Short: "Reancorá e atualiza uma série a partir de uma data",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			err = svc.AtualizarSerie(cmdContext(cmd), breno.AtualizacaoSerie{
				RecorrenciaID: upID, De: upDe, NovaData: upNova,
				Tipo: breno.TipoMovimentacao(upTipo), Valor: upValor,
				Descricao: upDescricao, TagIDs: upTags,
			})
			if err != nil {
				return err
			}
			return okMessage(cmd, d.g, "série atualizada")
		},
	}
	af := atualizar.Flags()
	af.StringVar(&upID, "id", "", "ID da recorrência (obrigatório)")
	af.StringVar(&upDe, "de", "", "âncora atual YYYY-MM-DD (obrigatório)")
	af.StringVar(&upNova, "nova-data", "", "nova âncora YYYY-MM-DD (obrigatório)")
	af.StringVar(&upTipo, "tipo", "", "tipo (obrigatório)")
	af.Float64Var(&upValor, "valor", 0, "novo valor (obrigatório)")
	af.StringVar(&upDescricao, "descricao", "", "nova descrição (obrigatório)")
	af.StringSliceVar(&upTags, "tags", nil, "IDs de tags")
	for _, r := range []string{"id", "de", "nova-data", "tipo", "valor", "descricao"} {
		_ = atualizar.MarkFlagRequired(r)
	}

	var (
		delID, delDe string
	)
	deletar := &cobra.Command{
		Use:   "deletar",
		Short: "Remove ocorrências de uma série a partir de uma data",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			if err := svc.DeletarSerieDaData(cmdContext(cmd), delID, delDe); err != nil {
				return err
			}
			return okMessage(cmd, d.g, "ocorrências removidas")
		},
	}
	deletar.Flags().StringVar(&delID, "id", "", "ID da recorrência (obrigatório)")
	deletar.Flags().StringVar(&delDe, "de", "", "data a partir da qual remover YYYY-MM-DD (obrigatório)")
	_ = deletar.MarkFlagRequired("id")
	_ = deletar.MarkFlagRequired("de")

	var delAllID string
	deletarTudo := &cobra.Command{
		Use:   "deletar-tudo",
		Short: "Remove todas as ocorrências de uma série",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			if err := svc.DeletarSerieToda(cmdContext(cmd), delAllID); err != nil {
				return err
			}
			return okMessage(cmd, d.g, "série removida")
		},
	}
	deletarTudo.Flags().StringVar(&delAllID, "id", "", "ID da recorrência (obrigatório)")
	_ = deletarTudo.MarkFlagRequired("id")

	group.AddCommand(atualizar, deletar, deletarTudo)
	return group
}
