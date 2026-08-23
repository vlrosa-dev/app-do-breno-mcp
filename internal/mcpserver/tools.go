package mcpserver

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

// Expenses

func ToolCreateExpenseOne(api *api.API) mcp.ToolHandlerFor[*structs.DespesasUnicaRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasUnicaRequest,
	) (*mcp.CallToolResult, any, error) {
		input.User = api.UserID

		tipoDespesa := strings.ToUpper(input.Tipo)
		opcoesDespesas := map[string]bool{
			"DIARIO":   true,
			"SAIDA":    true,
			"ECONOMIA": true,
			"ENTRADA":  true,
			"CARTAO":   true,
		}
		if !opcoesDespesas[tipoDespesa] {
			return textResult(
				"tipo invalido, valores esperados: Diario, Saida, Economia, Entrada ou Cartao",
			), nil, nil
		}
		id, err := api.CreateExpenseOne(ctx, input)
		if err != nil {
			return textResult("Unable to create transaction."), nil, nil
		}
		return textResult(id), nil, nil
	}
}

func ToolDeleteExpenseOne(api *api.API) mcp.ToolHandlerFor[*structs.DespesasUnicaDeleteRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasUnicaDeleteRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		if err := api.DeleteExpenseOne(ctx, input); err != nil {
			return textResult("Unable to delete transaction."), nil, nil
		}
		return textResult("Transaction deleted successfully."), nil, nil
	}
}

func ToolCreateExpensesRecurrency(api *api.API) mcp.ToolHandlerFor[*structs.DespesasRecorrenteRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasRecorrenteRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID

		_, err := time.Parse("2006-01-02", input.DataInicial)
		if err != nil {
			return textResult("error parsing date to format 2006-01-02"), nil, nil
		}

		if input.Infinita {
			input.NumeroRepeticoes = 24
			input.QuantidadeRepeticoes = 24
		}
		if input.NumeroRepeticoes <= 0 {
			return textResult(
				"numero de repeticoes deve ser maior que 0",
			), nil, nil
		}
		if input.QuantidadeRepeticoes <= 0 {
			return textResult(
				"quantidade de repeticoes deve ser maior que 0",
			), nil, nil
		}

		frequencia := strings.ToUpper(input.Frequencia)
		opcoesFrequencia := map[string]bool{
			"MENSALMENTE":  true,
			"SEMANALMENTE": true,
			"DIARIAMENTE":  true,
			"PARCELADO":    true,
		}
		if !opcoesFrequencia[frequencia] {
			return textResult(
				"tipo invalido, valores esperados: mensalmente, semanalmente, diariamente ou parcelado",
			), nil, nil
		}

		periodRecorrencia := strings.ToUpper(input.PeriodicidadeRecorrencia)
		opcoesPeriodicidadeRecorrencia := map[string]bool{
			"MENSALMENTE":  true,
			"SEMANALMENTE": true,
			"DIARIAMENTE":  true,
			"PARCELADO":    true,
		}
		if !opcoesPeriodicidadeRecorrencia[periodRecorrencia] {
			return textResult(
				"tipo invalido, valores esperados: mensalmente, semanalmente, diariamente ou parcelado",
			), nil, nil
		}

		tipo := strings.ToUpper(input.Tipo)
		opcoesDespesas := map[string]bool{
			"DIARIO":   true,
			"SAIDA":    true,
			"ECONOMIA": true,
			"ENTRADA":  true,
			"CARTAO":   true,
		}
		if !opcoesDespesas[tipo] {
			return textResult(
				"tipo invalido, valores esperados: Diario, Saida, Economia, Entrada ou Cartao",
			), nil, nil
		}

		if err := api.CreateExpensesRecurrency(ctx, *input); err != nil {
			return textResult(
				fmt.Sprintf("Unable to create recurring transaction, msg: %v", err.Error()),
			), nil, nil
		}
		return textResult("Recurring transaction created successfully."), nil, nil
	}
}

func ToolDeleteExpensesRecurrency(api *api.API) mcp.ToolHandlerFor[*structs.DespesasRecorrentesDeleteRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasRecorrentesDeleteRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		if err := api.DeleteExpensesRecurrency(ctx, *input); err != nil {
			return textResult(
				fmt.Sprintf("Unable to delete recurring transactions, msg: %v", err.Error()),
			), nil, nil
		}
		return textResult("Recurring transactions deleted successfully."), nil, nil
	}
}

func ToolDeleteExpensesRecurrencyAll(api *api.API) mcp.ToolHandlerFor[*structs.DespesasRecorrentesTodosDeleteRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasRecorrentesTodosDeleteRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		if err := api.DeleteExpensesRecurrencyAll(ctx, *input); err != nil {
			return textResult(
				fmt.Sprintf("Unable to delete all recurring transactions, msg: %v", err.Error()),
			), nil, nil
		}
		return textResult("All recurring transactions deleted successfully."), nil, nil
	}
}

func ToolListExpenses(api *api.API) mcp.ToolHandlerFor[*structs.DespesasListarRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasListarRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		if input.Select == "" {
			input.Select = "*"
		}

		dataInicial, err := time.Parse("2006-01-02", input.DataInicial)
		if err != nil {
			return textResult("error parsing date to format 2006-01-02"), nil, nil
		}
		dataFinal, err := time.Parse("2006-01-02", input.DataFinal)
		if err != nil {
			return textResult("error parsing date to format 2006-01-02"), nil, nil
		}

		if dataInicial.Equal(dataFinal) {
			input.DataFinal = dataFinal.AddDate(0, 0, 1).Format("2006-01-02")
		}

		despesas, err := api.ListExpenses(ctx, *input)
		if err != nil {
			return textResult(
				fmt.Sprintf("Unable to list transactions or no transactions found, msg: %v", err.Error()),
			), nil, nil
		}
		result, err := jsonResult(despesas)
		return result, listResult[structs.DespesasListarResponse]{Items: despesas}, err
	}
}

// Resumo

func ToolListResume(api *api.API) mcp.ToolHandlerFor[*structs.ResumoSaldosRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.ResumoSaldosRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserID = api.UserID
		resumo, err := api.ListResume(ctx, *input)
		if err != nil {
			return textResult(
				fmt.Sprintf("Unable to fetch resume, msg: %v", err.Error()),
			), nil, nil
		}
		result, err := jsonResult(resumo)
		return result, listResult[structs.ResumoSaldosResponse]{Items: resumo}, err
	}
}

// Tags

func ToolListTags(api *api.API) mcp.ToolHandlerFor[*structs.TagListRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TagListRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		tags, err := api.ListTags(ctx, *input)
		if err != nil {
			return textResult("Unable to list tags."), nil, nil
		}
		result, err := jsonResult(tags)
		return result, listResult[structs.TagListResponse]{Items: tags}, err
	}
}

// Cartao

func ToolListCards(api *api.API) mcp.ToolHandlerFor[*structs.ListCartaoRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.ListCartaoRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		if input.Select == "" {
			input.Select = "*"
		}
		cartoes, err := api.ListCards(ctx, *input)
		if err != nil {
			return textResult(
				fmt.Sprintf("Unable to list cards, msg: %v", err.Error()),
			), nil, nil
		}
		result, err := jsonResult(cartoes)
		return result, listResult[structs.ListCartaoResponse]{Items: cartoes}, err
	}
}

// Totais

func ToolListMonthlyTotals(api *api.API) mcp.ToolHandlerFor[*structs.TotaisMesListRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TotaisMesListRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserInput = api.UserID
		totais, err := api.ListResultstMonthly(ctx, *input)
		if err != nil {
			return textResult("Unable to fetch monthly totals."), nil, nil
		}
		result, err := jsonResult(totais)
		return result, totais, err
	}
}

func ToolListPerformanceYear(api *api.API) mcp.ToolHandlerFor[*structs.TotaisAnoPerformanceRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TotaisAnoPerformanceRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		performance, err := api.ListPerformeceYear(ctx, *input)
		if err != nil {
			return textResult("Unable to fetch yearly performance."), nil, nil
		}
		result, err := jsonResult(performance)
		return result, performance, err
	}
}

func ToolListLifeCostYear(api *api.API) mcp.ToolHandlerFor[*structs.TotaisAnoCustoDeVidaRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TotaisAnoCustoDeVidaRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		custoDeVida, err := api.ListLifeCostYear(ctx, *input)
		if err != nil {
			return textResult("Unable to fetch yearly life cost."), nil, nil
		}
		result, err := jsonResult(custoDeVida)
		return result, custoDeVida, err
	}
}

func ToolListSavingYear(api *api.API) mcp.ToolHandlerFor[*structs.TotaisAnoEconomizadoRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TotaisAnoEconomizadoRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		saving, err := api.ListSavingYear(ctx, *input)
		if err != nil {
			return textResult("Unable to fetch yearly savings."), nil, nil
		}
		result, err := jsonResult(saving)
		return result, saving, err
	}
}

func ToolListDailyAverageYear(api *api.API) mcp.ToolHandlerFor[*structs.TotaisDiarioMedioAnoRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.TotaisDiarioMedioAnoRequest,
	) (*mcp.CallToolResult, any, error) {
		input.UserId = api.UserID
		dailyAverage, err := api.ListDailyAverageYear(ctx, *input)
		if err != nil {
			return textResult("Unable to fetch yearly daily average."), nil, nil
		}
		result, err := jsonResult(dailyAverage)
		return result, dailyAverage, err
	}
}
