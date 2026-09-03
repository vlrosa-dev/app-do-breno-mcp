package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"app-do-breno-agent/internal/breno"

	"github.com/spf13/cobra"
)

// render escreve `data` no formato escolhido pela flag global --output.
func render(cmd *cobra.Command, g *globalFlags, data any) error {
	w := cmd.OutOrStdout()
	if g.output == "json" {
		return writeJSON(w, data)
	}
	return writeTable(w, data)
}

func writeJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeTable(w io.Writer, data any) error {
	tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
	defer tw.Flush()

	switch v := data.(type) {
	case []breno.Movimentacao:
		fmt.Fprintln(tw, "ID\tDATA\tTIPO\tVALOR\tDESCRIÇÃO")
		for _, m := range v {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", m.MovimentacaoID, m.Data, m.Tipo, money(m.Valor), m.Descricao)
		}
	case []breno.Saldo:
		fmt.Fprintln(tw, "MÊS\tENTRADAS\tSAÍDAS\tDIÁRIOS\tCARTÃO\tECONOMIAS\tSALDO")
		for _, s := range v {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				s.MesRef, money(s.TotalEntradas), money(s.TotalSaidas), money(s.TotalDiarios),
				money(s.TotalCartao), money(s.TotalEconomias), money(s.SaldoDia))
		}
	case []breno.Tag:
		fmt.Fprintln(tw, "TAG\tTOTAL")
		for _, t := range v {
			fmt.Fprintf(tw, "%s\t%s\n", t.Nome, money(t.TotalValor))
		}
	case []breno.CartaoCredito:
		fmt.Fprintln(tw, "NOME\tFECHAMENTO\tVENCIMENTO\tID")
		for _, c := range v {
			fmt.Fprintf(tw, "%s\t%d\t%d\t%s\n", c.Nome, c.DiaFechamento, c.DiaVencimento, c.Id)
		}
	case breno.TotaisMes:
		fmt.Fprintf(tw, "Performance\t%s\n", money(v.Performance))
		fmt.Fprintf(tw, "Custo de vida\t%s\n", money(v.CustoDeVida))
		fmt.Fprintf(tw, "Diário médio\t%s\n", money(v.DiarioMedio))
		fmt.Fprintf(tw, "Entradas\t%s\n", money(v.TotalEntradas))
		fmt.Fprintf(tw, "Saídas\t%s\n", money(v.TotalSaidas))
		fmt.Fprintf(tw, "Diários\t%s\n", money(v.TotalDiarios))
		fmt.Fprintf(tw, "Cartões\t%s\n", money(v.TotalCartoes))
		fmt.Fprintf(tw, "Economias\t%s\n", money(v.TotalEconomias))
	case string:
		fmt.Fprintln(tw, v)
	default:
		return writeJSON(w, data)
	}
	return nil
}

func money(v float64) string {
	return "R$ " + strconv.FormatFloat(v, 'f', 2, 64)
}

// okMessage imprime uma linha de sucesso simples (usado por comandos de escrita).
func okMessage(cmd *cobra.Command, g *globalFlags, msg string) error {
	if g.output == "json" {
		return writeJSON(cmd.OutOrStdout(), map[string]string{"status": "ok", "mensagem": msg})
	}
	fmt.Fprintln(cmd.OutOrStdout(), msg)
	return nil
}
