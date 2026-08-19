package structs

import "time"

type MovimentacaoRequest struct {
	Data      string  `json:"p_data" jsonschema:"required" description:"Data da movimentação, formato: YYYY-MM-DD"`
	Descricao string  `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	TagIds    []int   `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo      string  `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario' ou 'Entrada'"`
	User      string  `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor     float64 `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type ListarMovimentacoesRequest struct {
	// Select      string `params:"select" description:"Campos a serem selecionados na consulta, formato: *"`
	UserId      string `params:"user_id" description:"ID do usuário para filtrar as movimentações, formato: <user_id>"`
	Tipo        string `params:"tipo" description:"Tipo de movimentação para filtrar, valores possíveis: 'Diario' ou 'Entrada'"`
	DataInicial string `params:"data" description:"Data da movimentação para filtrar, formato: YYYY-MM-DD"`
	DataFinal   string `params:"data" description:"Data final da movimentação para filtrar, formato: YYYY-MM-DD"`
}

type ListarMovimentacoesResponse struct {
	MovimentacaoID       string        `json:"movimentacao_id"`
	CreatedAt            time.Time     `json:"created_at"`
	Valor                float64       `json:"valor"`
	Data                 string        `json:"data"`
	Descricao            string        `json:"descricao"`
	MovimentacaoOriginal interface{}   `json:"movimentacao_original"`
	Tipo                 string        `json:"tipo"`
	UserID               string        `json:"user_id"`
	RecorrenciaID        interface{}   `json:"recorrencia_id"`
	TagsJSON             []interface{} `json:"tags_json"`
	TagIds               []interface{} `json:"tag_ids"`
	TagNomes             []interface{} `json:"tag_nomes"`
}
