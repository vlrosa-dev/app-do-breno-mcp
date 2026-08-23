package mcpserver

import (
	"app-do-breno-agent/internal/api"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewMcpServer(api *api.API) (*mcp.Server, error) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:        "App do Breno MCP",
		Title:       "App do Breno",
		Description: "Servidor MCP com objetivo de lidar com as transações financeiras do APP Do Breno",
		Version:     "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_transacao_unica",
		Description: "Cria uma nova transação (movimentação) única",
	}, ToolCreateExpenseOne(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_transacao_unica",
		Description: "Deleta uma transação (movimentação) única",
	}, ToolDeleteExpenseOne(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_transacao_recorrente",
		Description: "Cria uma nova série de transações recorrentes",
	}, ToolCreateExpensesRecurrency(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_transacao_recorrente",
		Description: "Deleta as transações de uma série recorrente a partir de uma data",
	}, ToolDeleteExpensesRecurrency(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_transacao_recorrente_todas",
		Description: "Deleta todas as transações de uma série recorrente",
	}, ToolDeleteExpensesRecurrencyAll(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_transacoes",
		Description: "Lista as transações (movimentações) filtrando por tipo e período",
	}, ToolListExpenses(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_resumo_saldos",
		Description: "Lista o resumo diário de saldos de um ou mais meses",
	}, ToolListResume(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_tags",
		Description: "Lista as tags e seus totais em um mês",
	}, ToolListTags(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_totais_mes",
		Description: "Lista os totais e indicadores de um mês específico",
	}, ToolListMonthlyTotals(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_performance_ano",
		Description: "Lista a performance financeira mês a mês em um ano",
	}, ToolListPerformanceYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_custo_de_vida_ano",
		Description: "Lista o custo de vida mês a mês em um ano",
	}, ToolListLifeCostYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_economia_ano",
		Description: "Lista o total economizado mês a mês em um ano",
	}, ToolListSavingYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_diario_medio_ano",
		Description: "Lista o gasto diário médio mês a mês em um ano",
	}, ToolListDailyAverageYear(api))

	return server, nil
}
