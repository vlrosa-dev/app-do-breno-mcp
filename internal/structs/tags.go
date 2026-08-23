package structs

type TagListRequest struct {
	DataMes string `json:"p_data_mes" description:"Data no formato YYYY-MM-dd para filtrar as tags"`
	Search  string `json:"p_search" description:"Termo de busca para filtrar as tags"`
	UserId  string `json:"p_user_id" description:"ID do usuário para filtrar as tags"`
}

type TagListResponse struct {
	Id         string  `json:"tag_id"`
	Nome       string  `json:"tag_nome"`
	Cor        int     `json:"tag_cor"`
	TotalValor float64 `json:"total_valor"`
}
