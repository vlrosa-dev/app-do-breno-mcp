package structs

type TotaisMesListRequest struct {
	MesInput  string `json:"mes_input"`
	UserInput string `json:"user_input"`
}

type TotaisMesListResponse struct {
	Performance           float64 `json:"performance"`
	PerformancePrevisao   float64 `json:"performance_previsao"`
	PercentualEconomizado float64 `json:"percentual_economizado"`
	CustoDeVida           float64 `json:"custo_de_vida"`
	CustoDeVidaPrevisao   float64 `json:"custo_de_vida_previsao"`
	DiarioMedio           float64 `json:"diario_medio"`
	DiarioMedioAteHoje    float64 `json:"diario_medio_ate_hoje"`
	DiasNoMes             int     `json:"dias_no_mes"`
	DiasNoMesAteHoje      int     `json:"dias_no_mes_ate_hoje"`
	TotalEntradas         float64 `json:"total_entradas"`
	TotalSaidas           float64 `json:"total_saidas"`
	TotalDiarios          float64 `json:"total_diarios"`
	TotalEconomias        float64 `json:"total_economias"`
	TotalCartoes          float64 `json:"total_cartoes"`
	PrevisaoDiario        float64 `json:"previsao_diario"`
}

type TotaisAnoPerformanceRequest struct {
	Ano    string `json:"p_ano"`
	UserId string `json:"p_user_id"`
}

type TotaisAnoPerformaceResponse struct {
	Ano        int64 `json:"ano"`
	MesesNoAno int64 `json:"meses_no_ano"`
	Cabecalho  struct {
		PerformanceTotal float64 `json:"performance_total"`
	} `json:"cabecalho"`
	Meses []struct {
		MesNumero           int64   `json:"mes_numero"`
		MesNome             string  `json:"mes_nome"`
		Performance         float64 `json:"performance"`
		PerformancePrevisao float64 `json:"performance_previsao"`
	} `json:"meses"`
}

type TotaisAnoCustoDeVidaRequest struct {
	Ano    string `json:"p_ano"`
	UserId string `json:"p_user_id"`
}

type TotaisAnoCustoDeVidaResponse struct {
	Ano        int64 `json:"ano"`
	MesesNoAno int64 `json:"meses_no_ano"`
	Cabecalho  struct {
		CustoDeVidaTotal         float64 `json:"custo_de_vida_total"`
		CustoDeVidaTotalPrevisao float64 `json:"custo_de_vida_total_previsao"`
	} `json:"cabecalho"`
	Meses []struct {
		MesNumero           int64   `json:"mes_numero"`
		MesNome             string  `json:"mes_nome"`
		CustoDeVida         float64 `json:"custo_de_vida"`
		CustoDeVidaPrevisao float64 `json:"custo_de_vida_previsao"`
	} `json:"meses"`
}

type TotaisAnoEconomizadoRequest struct {
	Ano    string `json:"p_ano"`
	UserId string `json:"p_user_id"`
}

type TotaisAnoEconomizadoResponse struct {
	Ano        int64 `json:"ano"`
	MesesNoAno int64 `json:"meses_no_ano"`
	Cabecalho  struct {
		TotalEntradas         float64 `json:"total_entradas"`
		TotalEconomias        float64 `json:"total_economias"`
		PercentualEconomizado float64 `json:"percentual_economizado"`
	} `json:"cabecalho"`
	Meses []struct {
		MesNumero             int64   `json:"mes_numero"`
		MesNome               string  `json:"mes_nome"`
		TotalEntradas         float64 `json:"total_entradas"`
		TotalEconomias        float64 `json:"total_economias"`
		PercentualEconomizado float64 `json:"percentual_economizado"`
	} `json:"meses"`
}

type TotaisDiarioMedioAnoRequest struct {
	Ano    string `json:"p_ano"`
	UserId string `json:"p_user_id"`
}

type TotaisDiarioMedioAnoResponse struct {
	Ano        int64 `json:"ano"`
	MesesNoAno int64 `json:"meses_no_ano"`
	Cabecalho  struct {
		DiarioMedioAno float64 `json:"diario_medio_medio_ano"`
	} `json:"cabecalho"`
	Meses []struct {
		MesNumero    int64   `json:"mes_numero"`
		MesNome      string  `json:"mes_nome"`
		DiasNoMes    float64 `json:"dias_no_mes"`
		TotalDiarios float64 `json:"total_diarios"`
		DiarioMedio  float64 `json:"diario_medio"`
	} `json:"meses"`
}
