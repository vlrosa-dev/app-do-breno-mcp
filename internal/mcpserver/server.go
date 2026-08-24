package mcpserver

import (
	"app-do-breno-agent/internal/api"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BoolPtr(b bool) *bool {
	return &b
}

func NewMcpServer(logger *slog.Logger, api *api.API) (*mcp.Server, error) {
	implementation := &mcp.Implementation{
		Name:        "mcp-app-do-breno",
		Title:       "1.0.0",
		Description: "Servidor MCP com objetivo de lidar com as transações financeiras do APP Do Breno",
		Version:     "1.0.0",
	}

	options := &mcp.ServerOptions{
		Logger:    logger,
		KeepAlive: 5 * time.Minute,
		Instructions: `
		Você é um assistente financeiro e seu papel é fazer analises utilizando as ferramentas de listagem e realizar movimentações 
		financeiras no aplicativo utilizando as ferramentas de movimentacao, segue tipos de movimentações existentes:
		Tipos de movimentação:
			- Diario 	— gasto do dia a dia (débito)
			- Saida 	— saída fixa mensal
			- Cartao 	— despesa em cartão de crédito, pode ser parcelada ou unica
			- Entrada 	— receitas recebidas
			- Economia 	— valor economizado/poupado
		`,
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{},
		},
	}

	server := mcp.NewServer(implementation, options)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_movimentacao",
		Description: "Cria uma nova movimentação que pode ser do tipo (Entrada, Saida, Diario, Economia ou Cartão)",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleCreateTransactionTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atualizar_movimentacao",
		Description: "Atualiza uma movimentação que pode ser do tipo (Entrada, Saida, Diario, Economia ou Cartão)",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleUpdateTransactionTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao",
		Description: "Deleta movimentação a partir do seu ID",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleDeleteTransactionTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "criar_movimentacao_em_serie",
		Description: "Cria uma série de movimentações recorrentes que se repetem semanalmente, mensalmente, diariamente ou parcelado.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleCreateTransactionInSerieTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao_em_serie_da_data",
		Description: "Deleta as transações de uma série recorrente a partir de uma data e ID da movimentação recorrente",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleDeleteTransactionInSerieTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "deletar_todas_movimentacoes_em_serie",
		Description: "Deleta todas as transações recorrentes independente de data",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleDeleteAllTransactionInSerieTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_movimentacoes",
		Description: "Lista todas as movimentações filtrando por tipo (Entrada, Saida, Diario, Economia ou Cartao) e pelo período",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	}, HandleFetchTransactionTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_saldos",
		Description: "Lista o saldos para cada tipo a partir do mês solicitado",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchResumeBalanceTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_tags",
		Description: "Lista as tags e seus totais no mês",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchTagsTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_cartoes_de_creditos",
		Description: "Lista os cartões de creditos cadastrados com dia de vencimento e fechamento para calculo de fatura",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchCardsTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_totais_mes",
		Description: "Lista os totais e indicadores de um mês específico",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchTotalBalanceMonthlyTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_performance_ano",
		Description: "Lista a performance financeira mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchPerformanceYearTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_custo_de_vida_ano",
		Description: "Lista o custo de vida mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchLifeCostYearTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_economia_ano",
		Description: "Lista o total economizado mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchSavingYearTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_diario_medio_ano",
		Description: "Lista o gasto diário médio mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchAverageCostDailyTool(api))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listar_movimentacao_por_id",
		Description: "Obter movimentação pelo ID para atualização dos dados",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	}, HandleFetchTransactionByIdTool(api))

	return server, nil
}
