package structs

type ListCartaoRequest struct {
	Order  string
	Select string
	UserId string
}

type ListCartaoResponse struct {
	Id            string `json:"id" jsonschema:"required" description:"Id do cartão UUID"`
	CreatedAt     string `json:"created_at" jsonschema:"required" description:"Data de criação do cartão"`
	Nome          string `json:"nome" description:"Nome do Cartão"`
	DiaFechamento int    `json:"dia_fechamento" jsonschema:"required" description:"Dia do fechamento da fatura"`
	DiaVencimento int    `json:"dia_vencimento" jsonschema:"required" description:"Dia do vencimento da fatura"`
	UserId        string `json:"user_id" jsonschema:"required" description:"ID do usuario"`
}
