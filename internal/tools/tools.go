package tools

import (
	"app-do-breno-agent/internal/api"
	"app-do-breno-agent/internal/structs"
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ToolCreateTransaction(api *api.API) mcp.ToolHandlerFor[*structs.DespesasUnicaRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasUnicaRequest,
	) (*mcp.CallToolResult, any, error) {
		id_transacao, err := api.CreateNewTransaction(ctx, input)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Unable to create transaction or no transactions found."},
				},
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: id_transacao},
			},
		}, nil, nil
	}
}

func ToolFetchTransactions(api *api.API) mcp.ToolHandlerFor[*structs.DespesasListarRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasListarRequest,
	) (*mcp.CallToolResult, any, error) {
		movimentacoes, err := api.FetchTransactions(ctx, input)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Unable to list transactions or no transactions found."},
				},
			}, nil, nil
		}

		movimentacoes_json, err := json.Marshal(movimentacoes)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error occurred while marshaling transactions."},
				},
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(movimentacoes_json)},
			},
		}, movimentacoes, nil
	}
}

func ToolDeleteTransaction(api *api.API) mcp.ToolHandlerFor[*structs.DespesasUnicaDeleteRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.DespesasUnicaDeleteRequest,
	) (*mcp.CallToolResult, any, error) {
		err := api.DeleteTransaction(ctx, "")
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Unable to delete transaction."},
				},
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction deleted successfully."},
			},
		}, nil, nil
	}
}
