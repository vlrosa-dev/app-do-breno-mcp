package cli

import (
	"context"

	"app-do-breno-agent/internal/breno"

	"github.com/spf13/cobra"
)

// newConsultaCmds cria os comandos de consulta de nível superior.
func newConsultaCmds(d *deps) []*cobra.Command {
	var (
		saldosMes   string
		saldosMeses int
	)
	saldos := &cobra.Command{
		Use:   "saldos",
		Short: "Saldos por tipo a partir de um mês",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			res, err := svc.ListarSaldos(cmdContext(cmd), saldosMes, saldosMeses)
			if err != nil {
				return err
			}
			return render(cmd, d.g, res)
		},
	}
	saldos.Flags().StringVar(&saldosMes, "mes", "", "mês YYYY-MM-DD (obrigatório)")
	saldos.Flags().IntVar(&saldosMeses, "meses", 1, "quantidade de meses")
	_ = saldos.MarkFlagRequired("mes")

	var (
		tagsMes   string
		tagsBusca string
	)
	tags := &cobra.Command{
		Use:   "tags",
		Short: "Tags e seus totais no mês",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			res, err := svc.ListarTags(cmdContext(cmd), tagsMes, tagsBusca)
			if err != nil {
				return err
			}
			return render(cmd, d.g, res)
		},
	}
	tags.Flags().StringVar(&tagsMes, "mes", "", "mês YYYY-MM-DD (obrigatório)")
	tags.Flags().StringVar(&tagsBusca, "busca", "", "termo de busca")
	_ = tags.MarkFlagRequired("mes")

	cartoes := &cobra.Command{
		Use:   "cartoes",
		Short: "Cartões de crédito cadastrados",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			res, err := svc.ListarCartoes(cmdContext(cmd))
			if err != nil {
				return err
			}
			return render(cmd, d.g, res)
		},
	}

	return []*cobra.Command{saldos, tags, cartoes}
}

// newTotaisGroupCmd cria o grupo `totais`.
func newTotaisGroupCmd(d *deps) *cobra.Command {
	group := &cobra.Command{Use: "totais", Short: "Indicadores agregados"}

	var mesInput string
	mes := &cobra.Command{
		Use:   "mes",
		Short: "Indicadores de um mês",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := d.service()
			if err != nil {
				return err
			}
			res, err := svc.TotaisMes(cmdContext(cmd), mesInput)
			if err != nil {
				return err
			}
			return render(cmd, d.g, res)
		},
	}
	mes.Flags().StringVar(&mesInput, "mes", "", "mês YYYY-MM-DD ou YYYY-MM (obrigatório)")
	_ = mes.MarkFlagRequired("mes")

	anoCmd := func(use, short string, run func(svc breno.Financas, ctx context.Context, ano string) (any, error)) *cobra.Command {
		var ano string
		c := &cobra.Command{
			Use:   use,
			Short: short,
			RunE: func(cmd *cobra.Command, _ []string) error {
				svc, err := d.service()
				if err != nil {
					return err
				}
				res, err := run(svc, cmdContext(cmd), ano)
				if err != nil {
					return err
				}
				return render(cmd, d.g, res)
			},
		}
		c.Flags().StringVar(&ano, "ano", "", "ano (ex.: 2026) (obrigatório)")
		_ = c.MarkFlagRequired("ano")
		return c
	}

	group.AddCommand(
		mes,
		anoCmd("performance", "Performance mês a mês no ano", func(svc breno.Financas, ctx context.Context, ano string) (any, error) {
			return svc.PerformanceAno(ctx, ano)
		}),
		anoCmd("custo-vida", "Custo de vida mês a mês no ano", func(svc breno.Financas, ctx context.Context, ano string) (any, error) {
			return svc.CustoDeVidaAno(ctx, ano)
		}),
		anoCmd("economia", "Total economizado mês a mês no ano", func(svc breno.Financas, ctx context.Context, ano string) (any, error) {
			return svc.EconomiaAno(ctx, ano)
		}),
		anoCmd("diario-medio", "Gasto diário médio mês a mês no ano", func(svc breno.Financas, ctx context.Context, ano string) (any, error) {
			return svc.DiarioMedioAno(ctx, ano)
		}),
	)
	return group
}
