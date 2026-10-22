package mcpserver

import (
	"github.com/vlrosa-dev/app-do-breno-mcp/internal/breno"
	"log/slog"
	"reflect"
	"slices"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BoolPtr(b bool) *bool {
	return &b
}

// camposOpcionais são propriedades que o servidor preenche sozinho (o id do
// usuário vem do login) ou ignora. Continuam aceitas na entrada, para não
// quebrar agentes que já as enviam, mas deixam de ser obrigatórias.
var camposOpcionais = []string{
	"p_user", "p_user_id", "user_id", "user_input", "UserId",
	"Order", "Select", "Limit", "select",
}

// inputSchema gera o schema de entrada a partir do struct, como o SDK faria
// (derreferenciando ponteiros), e remove camposOpcionais da lista de obrigatórios.
func inputSchema[In any]() (*jsonschema.Schema, error) {
	rt := reflect.TypeFor[In]()
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	schema, err := jsonschema.ForType(rt, &jsonschema.ForOptions{})
	if err != nil {
		return nil, err
	}
	schema.Required = slices.DeleteFunc(schema.Required, func(name string) bool {
		return slices.Contains(camposOpcionais, name)
	})
	return schema, nil
}

// addTool registra a tool com o schema de entrada ajustado por inputSchema.
func addTool[In any](server *mcp.Server, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, any]) error {
	schema, err := inputSchema[In]()
	if err != nil {
		return err
	}
	tool.InputSchema = schema
	mcp.AddTool(server, tool, handler)
	return nil
}

func NewMcpServer(logger *slog.Logger, svc breno.Financas) (*mcp.Server, error) {
	implementation := &mcp.Implementation{
		Name:        "app-do-breno-mcp",
		Title:       "App do Breno MCP",
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

	if err := addTool(server, &mcp.Tool{
		Name:        "criar_movimentacao",
		Description: "Cria uma nova movimentação que pode ser do tipo (Entrada, Saida, Diario, Economia ou Cartão)",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleCreateTransactionTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "atualizar_movimentacao",
		Description: "Atualiza uma movimentação que pode ser do tipo (Entrada, Saida, Diario, Economia ou Cartão)",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleUpdateTransactionTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao",
		Description: "Deleta movimentação a partir do seu ID",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleDeleteTransactionTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "criar_movimentacao_em_serie",
		Description: "Cria uma série de movimentações recorrentes que se repetem semanalmente, mensalmente, diariamente ou parcelado.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleCreateTransactionInSerieTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "atualizar_movimentacao_em_serie",
		Description: "Atualiza os dados de uma transação recorrente, acione a tool 'listar_movimentacao_por_id' antes de atualizar.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleUpdateTransactionInSerieTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "deletar_movimentacao_em_serie_da_data",
		Description: "Deleta as transações de uma série recorrente a partir de uma data e ID da movimentação recorrente",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleDeleteTransactionInSerieTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "deletar_todas_movimentacoes_em_serie",
		Description: "Deleta todas as transações recorrentes independente de data",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleDeleteAllTransactionInSerieTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_movimentacoes",
		Description: "Lista todas as movimentações filtrando por tipo (Entrada, Saida, Diario, Economia ou Cartao) e pelo período",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    false,
		},
	},
		HandleFetchTransactionTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_saldos",
		Description: "Lista o saldos para cada tipo a partir do mês solicitado",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchResumeBalanceTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_tags",
		Description: "Lista as tags e seus totais no mês",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchTagsTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_cartoes_de_creditos",
		Description: "Lista os cartões de creditos cadastrados com dia de vencimento e fechamento para calculo de fatura",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchCardsTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_totais_mes",
		Description: "Lista os totais e indicadores de um mês específico",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchTotalBalanceMonthlyTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_performance_ano",
		Description: "Lista a performance financeira mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchPerformanceYearTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_custo_de_vida_ano",
		Description: "Lista o custo de vida mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchLifeCostYearTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_economia_ano",
		Description: "Lista o total economizado mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchSavingYearTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_diario_medio_ano",
		Description: "Lista o gasto diário médio mês a mês em um ano",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchAverageCostDailyTool(svc)); err != nil {
		return nil, err
	}

	if err := addTool(server, &mcp.Tool{
		Name:        "listar_movimentacao_por_id",
		Description: "Obter movimentação pelo ID para atualização dos dados",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: BoolPtr(true),
			IdempotentHint:  false,
			OpenWorldHint:   BoolPtr(false),
			ReadOnlyHint:    true,
		},
	},
		HandleFetchTransactionByIdTool(svc)); err != nil {
		return nil, err
	}

	return server, nil
}
