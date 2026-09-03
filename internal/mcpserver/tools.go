package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"app-do-breno-agent/internal/breno"
	"app-do-breno-agent/internal/structs"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Os structs de entrada das tools continuam sendo os de internal/structs para
// preservar exatamente o schema MCP já publicado. As regras de negócio agora
// vivem em internal/breno — os handlers apenas mapeiam entrada -> domínio.

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return textResult("Error occurred while marshaling response."), err
	}
	return textResult(string(data)), nil
}

// listResult wraps a slice so it can be used as MCP structuredContent, which
// the protocol requires to be a JSON object — a bare JSON array fails validation.
type listResult[T any] struct {
	Items []T `json:"items"`
}

func errText(prefix string, err error) *mcp.CallToolResult {
	return textResult(fmt.Sprintf("%s: %v", prefix, err))
}

// Transactions

func HandleCreateTransactionTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionRequest) (*mcp.CallToolResult, any, error) {
		id, err := svc.CriarMovimentacao(ctx, breno.NovaMovimentacao{
			Tipo:      breno.TipoMovimentacao(input.Tipo),
			Valor:     input.Valor,
			Descricao: input.Descricao,
			Data:      input.Data,
			TagIDs:    input.TagIds,
		})
		if err != nil {
			return errText("Unable to create transaction", err), nil, nil
		}
		return textResult(id), nil, nil
	}
}

func HandleUpdateTransactionTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionUpdateRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionUpdateRequest) (*mcp.CallToolResult, any, error) {
		err := svc.AtualizarMovimentacao(ctx, breno.AtualizacaoMovimentacao{
			ID:        input.MovId,
			Tipo:      breno.TipoMovimentacao(input.Tipo),
			Valor:     input.Valor,
			Descricao: input.Descricao,
			Data:      input.Data,
			TagIDs:    input.TagIds,
		})
		if err != nil {
			return errText("Unable to update transaction", err), nil, nil
		}
		return textResult("transaction updated successfully"), nil, nil
	}
}

func HandleDeleteTransactionTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionDeleteRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionDeleteRequest) (*mcp.CallToolResult, any, error) {
		if err := svc.DeletarMovimentacao(ctx, input.Id); err != nil {
			return errText("Unable to delete transaction", err), nil, nil
		}
		return textResult("Transaction deleted successfully."), nil, nil
	}
}

func HandleCreateTransactionInSerieTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionInSerieRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionInSerieRequest) (*mcp.CallToolResult, any, error) {
		_, err := svc.CriarMovimentacao(ctx, breno.NovaMovimentacao{
			Tipo:      breno.TipoMovimentacao(input.Tipo),
			Valor:     input.Valor,
			Descricao: input.Descricao,
			Data:      input.DataInicial,
			TagIDs:    input.TagIds,
			Recorrencia: &breno.Recorrencia{
				Frequencia: breno.Frequencia(input.Frequencia),
				Repeticoes: input.NumeroRepeticoes,
				Infinita:   input.Infinita,
			},
		})
		if err != nil {
			return errText("Unable to create recurring transaction", err), nil, nil
		}
		return textResult("Recurring transaction created successfully."), nil, nil
	}
}

func HandleUpdateTransactionInSerieTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionInSerieUpdateRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionInSerieUpdateRequest) (*mcp.CallToolResult, any, error) {
		err := svc.AtualizarSerie(ctx, breno.AtualizacaoSerie{
			RecorrenciaID: input.Recorrencia,
			De:            input.From,
			NovaData:      input.NewDate,
			Tipo:          breno.TipoMovimentacao(input.Tipo),
			Valor:         input.Valor,
			Descricao:     input.Descricao,
			TagIDs:        input.TagIds,
		})
		if err != nil {
			return errText("Unable to update recurring transaction", err), nil, nil
		}
		return textResult("transaction updated successfully"), nil, nil
	}
}

func HandleDeleteTransactionInSerieTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionInSerieDeleteRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionInSerieDeleteRequest) (*mcp.CallToolResult, any, error) {
		if err := svc.DeletarSerieDaData(ctx, input.RecorrenciaID, input.From); err != nil {
			return errText("Unable to delete recurring transactions", err), nil, nil
		}
		return textResult("Recurring transactions deleted successfully."), nil, nil
	}
}

func HandleDeleteAllTransactionInSerieTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionInSerieDeleteAllRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionInSerieDeleteAllRequest) (*mcp.CallToolResult, any, error) {
		if err := svc.DeletarSerieToda(ctx, input.RecorrenciaID); err != nil {
			return errText("Unable to delete all recurring transactions", err), nil, nil
		}
		return textResult("All recurring transactions deleted successfully."), nil, nil
	}
}

func HandleFetchTransactionTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionListRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionListRequest) (*mcp.CallToolResult, any, error) {
		movs, err := svc.ListarMovimentacoes(ctx, breno.FiltroMovimentacoes{
			Tipo:   breno.TipoMovimentacao(input.Tipo),
			De:     input.DataInicial,
			Ate:    input.DataFinal,
			Select: input.Select,
		})
		if err != nil {
			return errText("Unable to list transactions", err), nil, nil
		}
		result, err := jsonResult(movs)
		return result, listResult[breno.Movimentacao]{Items: movs}, err
	}
}

func HandleFetchTransactionByIdTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TransactionGetRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TransactionGetRequest) (*mcp.CallToolResult, any, error) {
		movs, err := svc.ObterMovimentacao(ctx, input.MovimentacaoId)
		if err != nil {
			return errText("Unable to fetch transaction by id", err), nil, nil
		}
		result, err := jsonResult(movs)
		return result, listResult[breno.Movimentacao]{Items: movs}, err
	}
}

// Resumo

func HandleFetchResumeBalanceTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.ResumoSaldosRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.ResumoSaldosRequest) (*mcp.CallToolResult, any, error) {
		saldos, err := svc.ListarSaldos(ctx, input.DataMes, int(input.QuantidadeMeses))
		if err != nil {
			return errText("Unable to fetch resume", err), nil, nil
		}
		result, err := jsonResult(saldos)
		return result, listResult[breno.Saldo]{Items: saldos}, err
	}
}

// Tags

func HandleFetchTagsTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TagListRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TagListRequest) (*mcp.CallToolResult, any, error) {
		tags, err := svc.ListarTags(ctx, input.DataMes, input.Search)
		if err != nil {
			return errText("Unable to list tags", err), nil, nil
		}
		result, err := jsonResult(tags)
		return result, listResult[breno.Tag]{Items: tags}, err
	}
}

// Cartao

func HandleFetchCardsTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.ListCartaoRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ *structs.ListCartaoRequest) (*mcp.CallToolResult, any, error) {
		cards, err := svc.ListarCartoes(ctx)
		if err != nil {
			return errText("Unable to list cards", err), nil, nil
		}
		result, err := jsonResult(cards)
		return result, listResult[breno.CartaoCredito]{Items: cards}, err
	}
}

// Totais

func HandleFetchTotalBalanceMonthlyTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TotaisMesListRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TotaisMesListRequest) (*mcp.CallToolResult, any, error) {
		totais, err := svc.TotaisMes(ctx, input.MesInput)
		if err != nil {
			return errText("Unable to fetch monthly totals", err), nil, nil
		}
		result, err := jsonResult(totais)
		return result, totais, err
	}
}

func HandleFetchPerformanceYearTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TotaisAnoPerformanceRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TotaisAnoPerformanceRequest) (*mcp.CallToolResult, any, error) {
		perf, err := svc.PerformanceAno(ctx, input.Ano)
		if err != nil {
			return errText("Unable to fetch yearly performance", err), nil, nil
		}
		result, err := jsonResult(perf)
		return result, perf, err
	}
}

func HandleFetchLifeCostYearTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TotaisAnoCustoDeVidaRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TotaisAnoCustoDeVidaRequest) (*mcp.CallToolResult, any, error) {
		resp, err := svc.CustoDeVidaAno(ctx, input.Ano)
		if err != nil {
			return errText("Unable to fetch yearly life cost", err), nil, nil
		}
		result, err := jsonResult(resp)
		return result, resp, err
	}
}

func HandleFetchSavingYearTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TotaisAnoEconomizadoRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TotaisAnoEconomizadoRequest) (*mcp.CallToolResult, any, error) {
		resp, err := svc.EconomiaAno(ctx, input.Ano)
		if err != nil {
			return errText("Unable to fetch yearly savings", err), nil, nil
		}
		result, err := jsonResult(resp)
		return result, resp, err
	}
}

func HandleFetchAverageCostDailyTool(svc breno.Financas) mcp.ToolHandlerFor[*structs.TotaisDiarioMedioAnoRequest, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *structs.TotaisDiarioMedioAnoRequest) (*mcp.CallToolResult, any, error) {
		resp, err := svc.DiarioMedioAno(ctx, input.Ano)
		if err != nil {
			return errText("Unable to fetch yearly daily average", err), nil, nil
		}
		result, err := jsonResult(resp)
		return result, resp, err
	}
}
