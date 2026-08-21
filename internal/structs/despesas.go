package structs

import "time"

// Expenses One
type DespesasUnicaRequest struct {
	Data      string  `json:"p_data" jsonschema:"required" description:"Data da movimentação, formato: YYYY-MM-DD"`
	Descricao string  `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	TagIds    []int   `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo      string  `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario' ou 'Entrada'"`
	User      string  `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor     float64 `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type DespesasUnicaDeleteRequest struct {
	Id     string `params:"id" description:"ID da movimentação a ser deletada"`
	UserId string `params:"user_id" description:"ID do usuário que está realizando a movimentação"`
}

// Expenses Recurrency
type DespesasRecorrenteRequest struct {
	DataInicial              string  `json:"p_data_inicial" jsonschema:"required" description:"Data inicial da recorrência, formato: YYYY-MM-DDTHH:MM:SSZ"`
	Descricao                string  `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	Frequencia               string  `json:"p_frequencia" jsonschema:"required" description:"Frequência da recorrência, valores possíveis: 'mensalmente', 'semanalmente', 'diariamente', 'parcelado'"`
	Infinita                 bool    `json:"p_infinita" jsonschema:"required" description:"Indica se a recorrência é infinita"`
	NumeroRepeticoes         int     `json:"p_numero_repeticoes" description:"Número de repetições da recorrência, necessário se a recorrência não for infinita"`
	PeriodicidadeRecorrencia string  `json:"p_periodicidade_recorrencia" description:"Periodicidade da recorrência, necessário se a recorrência não for infinita, valores possíveis: 'Mensalmente', 'Semanalmente', 'Diariamente'"`
	QuantidadeRepeticoes     int     `json:"p_quantidade_repeticoes" description:"Quantidade de repetições da recorrência, necessário se a recorrência não for infinita"`
	TagIds                   []int   `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo                     string  `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario', 'Entrada', 'Economia', 'Cartao'"`
	UserId                   string  `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor                    float64 `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type DespesasRecorrentesDeleteRequest struct {
	From          string `params:"p_from" description:"Data a partir da qual as movimentações recorrentes devem ser deletadas, formato: YYYY-MM-DD"`
	RecorrenciaID string `params:"p_recorrencia" description:"ID da recorrência a ser deletada"`
	UserId        string `params:"p_user" description:"ID do usuário que está realizando a movimentação"`
}

type DespesasRecorrentesTodosDeleteRequest struct {
	RecorrenciaID string `params:"p_recorrencia" description:"ID da recorrência a ser deletada"`
	UserId        string `params:"p_user" description:"ID do usuário que está realizando a movimentação"`
}

// LIST
type DespesasListarRequest struct {
	Select      string `params:"select" description:"Campos a serem selecionados na consulta, formato: *"`
	UserId      string `params:"user_id" description:"ID do usuário para filtrar as movimentações, formato: <user_id>"`
	Tipo        string `params:"tipo" description:"Tipo de movimentação para filtrar, valores possíveis: 'Diario' ou 'Entrada'"`
	DataInicial string `params:"data" description:"Data da movimentação para filtrar, formato: YYYY-MM-DD"`
	DataFinal   string `params:"data" description:"Data final da movimentação para filtrar, formato: YYYY-MM-DD"`
}

type DespesasListarResponse struct {
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
