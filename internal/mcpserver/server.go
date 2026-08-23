package mcpserver

import (
	"app-do-breno-agent/internal/api"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BoolPtr(b bool) *bool {
	return &b
}

func NewMcpServer(api *api.API) (*mcp.Server, error) {

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:        "mcp-app-do-breno",
			Title:       "1.0.0",
			Description: "Servidor MCP com objetivo de lidar com as transações financeiras do APP Do Breno",
			Version:     "1.0.0",
		}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_movimentacao",
		Description: "Cria uma nova movimentação que pode ser do tipo (Entrada, Saida, Diario, Economia ou Cartão)",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolCreateExpenseOne(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao",
		Description: "Deleta movimentação a partir do seu ID",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolDeleteExpenseOne(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_movimentacao_em_serie",
		Description: "Cria uma série de movimentações recorrentes que se repetem semanalmente, mensalmente, diariamente ou parcelado.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolCreateExpensesRecurrency(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao_em_serie_da_data",
		Description: "Deleta as transações de uma série recorrente a partir de uma data e ID da movimentação recorrente",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolDeleteExpensesRecurrency(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_todas_movimentacoes_em_serie",
		Description: "Deleta todas as transações recorrentes independente de data",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolDeleteExpensesRecurrencyAll(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_movimentacoes",
		Description: "Lista todas as movimentações filtrando por tipo (Entrada, Saida, Diario, Economia ou Cartao) e pelo período",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, ToolListExpenses(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_saldos",
		Description: "Lista o saldos para cada tipo a partir do mês solicitado",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListResume(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_tags",
		Description: "Lista as tags e seus totais no mês",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListTags(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_cartoes_de_creditos",
		Description: "Lista os cartões de creditos cadastrados com dia de vencimento e fechamento para calculo de fatura",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListCards(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_totais_mes",
		Description: "Lista os totais e indicadores de um mês específico",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListMonthlyTotals(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_performance_ano",
		Description: "Lista a performance financeira mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListPerformanceYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_custo_de_vida_ano",
		Description: "Lista o custo de vida mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListLifeCostYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_economia_ano",
		Description: "Lista o total economizado mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListSavingYear(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_diario_medio_ano",
		Description: "Lista o gasto diário médio mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, ToolListDailyAverageYear(api))

	return server, nil
}
