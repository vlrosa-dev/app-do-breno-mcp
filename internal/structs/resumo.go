package structs

type ResumoSaldosRequest struct {
	DataMes         string `json:"p_data_mes"`
	QuantidadeMeses int64  `json:"p_quantidade_meses"`
	UserID          string `json:"p_user_id"`
}

type ResumoSaldosResponse struct {
	DataISO        string  `json:"data_iso"`
	MesRef         string  `json:"mes_ref"`
	MesOrdem       string  `json:"mes_ordem"`
	TotalEntradas  float64 `json:"total_entradas"`
	TotalSaidas    float64 `json:"total_saidas"`
	TotalDiarios   float64 `json:"total_diarios"`
	TotalEconomias float64 `json:"total_economias"`
	TotalCartao    float64 `json:"total_cartao"`
	SaldoDia       float64 `json:"saldo_dia"`
	Checkin        bool    `json:"checkin"`
}
