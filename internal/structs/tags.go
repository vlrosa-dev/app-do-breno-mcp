package structs

type TagListRequest struct {
	DataMes string `params:"p_data_mes" description:"Data no formato YYYY-MM-dd para filtrar as tags"`
	Search  string `params:"p_search" description:"Termo de busca para filtrar as tags"`
	UserId  string `params:"p_user_id" description:"ID do usuário para filtrar as tags"`
}

type TagListResponse struct {
	Id         int     `json:"tag_id"`
	Nome       string  `json:"tag_nome"`
	Cor        string  `json:"tag_cor"`
	TotalValor float32 `json:"total_valor"`
}
