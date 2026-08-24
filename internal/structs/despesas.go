package structs

import "time"

type TransactionRequest struct {
	Data      string   `json:"p_data" jsonschema:"required" description:"Data da movimentação, formato: YYYY-MM-DD"`
	Descricao string   `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	TagIds    []string `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo      string   `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario' ou 'Entrada'"`
	User      string   `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor     float64  `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type TransactionUpdateRequest struct {
	Data      string   `json:"p_data" jsonschema:"required" description:"Data da movimentação, formato: YYYY-MM-DD"`
	Descricao string   `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	MovId     string   `json:"p_mov" jsonschema:"required" description:"Descrição da movimentação"`
	TagIds    []string `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo      string   `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario', 'Entrada', 'Cartao' ou 'Saida'"`
	User      string   `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor     float64  `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type TransactionGetRequest struct {
	Limit          int
	MovimentacaoId string
	Select         string
	UserId         string
}

type TransactionDeleteRequest struct {
	Id     string `json:"id" description:"ID da movimentação a ser deletada"`
	UserId string `json:"user_id" description:"ID do usuário que está realizando a movimentação"`
}

type TransactionInSerieRequest struct {
	DataInicial              string   `json:"p_data_inicial" jsonschema:"required" description:"Data inicial da recorrência, formato: YYYY-MM-DD"`
	Descricao                string   `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	Frequencia               string   `json:"p_frequencia" jsonschema:"required" description:"Frequência da recorrência, valores possíveis: 'mensalmente', 'semanalmente', 'diariamente', 'parcelado'"`
	Infinita                 bool     `json:"p_infinita" jsonschema:"required" description:"Indica se a recorrência é infinita"`
	NumeroRepeticoes         int      `json:"p_numero_repeticoes" description:"Número de repetições da recorrência, necessário se a recorrência não for infinita"`
	PeriodicidadeRecorrencia string   `json:"p_periodicidade_recorrencia" description:"Periodicidade da recorrência, necessário se a recorrência não for infinita, valores possíveis: 'Mensalmente', 'Semanalmente', 'Diariamente'"`
	QuantidadeRepeticoes     int      `json:"p_quantidade_repeticoes" description:"Quantidade de repetições da recorrência, necessário se a recorrência não for infinita"`
	TagIds                   []string `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo                     string   `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario', 'Entrada', 'Economia', 'Cartao'"`
	UserId                   string   `json:"p_user_id" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor                    float64  `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type TransactionInSerieUpdateRequest struct {
	DateMode    string   `json:"p_data_inicial" jsonschema:"required"`
	Descricao   string   `json:"p_descricao" jsonschema:"required" description:"Descrição da movimentação"`
	From        string   `json:"p_from" jsonschema:"required" description:"Data inicial da recorrência, formato: YYYY-MM-DD"`
	NewDate     string   `json:"p_new_date" jsonschema:"required" description:"Data inicial da recorrência, formato: YYYY-MM-DD"`
	Recorrencia string   `json:"p_recorrencia" jsonschema:"required" description:"ID da reccorencia"`
	TagIds      []string `json:"p_tag_ids" description:"IDs das tags associadas à movimentação"`
	Tipo        string   `json:"p_tipo" jsonschema:"required" description:"Tipo da movimentação, valores possíveis: 'Diario', 'Entrada', 'Economia', 'Cartao'"`
	User        string   `json:"p_user" jsonschema:"required" description:"Usuário que está realizando a movimentação"`
	Valor       float64  `json:"p_valor" jsonschema:"required" description:"Valor da movimentação"`
}

type TransactionInSerieDeleteRequest struct {
	From          string `json:"p_from" description:"Data a partir da qual as movimentações recorrentes devem ser deletadas, formato: YYYY-MM-DD"`
	RecorrenciaID string `json:"p_recorrencia" description:"ID da recorrência a ser deletada"`
	UserId        string `json:"p_user" description:"ID do usuário que está realizando a movimentação"`
}

type TransactionInSerieDeleteAllRequest struct {
	RecorrenciaID string `json:"p_recorrencia" description:"ID da recorrência a ser deletada"`
	UserId        string `json:"p_user" description:"ID do usuário que está realizando a movimentação"`
}

type TransactionListRequest struct {
	Select      string `json:"select" description:"Campos a serem selecionados na consulta, default: *"`
	UserId      string `json:"user_id" description:"ID do usuário para filtrar as movimentações, formato: <user_id>"`
	Tipo        string `json:"tipo" description:"Tipo de movimentação para filtrar, valores possíveis: 'Diario' ou 'Entrada'"`
	DataInicial string `json:"data_inicial" description:"Data da movimentação para filtrar, formato: YYYY-MM-DD"`
	DataFinal   string `json:"data_final" description:"Data final da movimentação para filtrar, formato: YYYY-MM-DD"`
}

type TransactionListResponse struct {
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
