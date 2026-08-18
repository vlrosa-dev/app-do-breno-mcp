package structs

type MovimentacaoRequest struct {
	Data      string  `json:"p_data"`
	Descricao string  `json:"p_descricao"`
	TagIds    []int   `json:"p_tag_ids"`
	Tipo      string  `json:"p_tipo"`
	User      string  `json:"p_user"`
	Valor     float64 `json:"p_valor"`
}

type MovimentacaoResponse struct {
	Id string `json:"id"`
}
