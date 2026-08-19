package tools

import (
	"app-do-breno-agent/internal/structs"
	"app-do-breno-agent/internal/www"
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ToolMcpMovimentacao(api *www.API) mcp.ToolHandlerFor[*structs.MovimentacaoRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.MovimentacaoRequest,
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

func ToolMcpListarMovimentacoes(api *www.API) mcp.ToolHandlerFor[*structs.ListarMovimentacoesRequest, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input *structs.ListarMovimentacoesRequest,
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
