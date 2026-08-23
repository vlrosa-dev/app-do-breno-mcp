package api

import (
	"app-do-breno-agent/internal/structs"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (api *API) ListResultstMonthly(ctx context.Context, obj structs.TotaisMesListRequest) (structs.TotaisMesListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_totais_mes")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.TotaisMesListResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return structs.TotaisMesListResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.TotaisMesListResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return structs.TotaisMesListResponse{}, err
	}

	if res == nil {
		return structs.TotaisMesListResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.TotaisMesListResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return structs.TotaisMesListResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse structs.TotaisMesListResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.TotaisMesListResponse{}, err
	}

	return objResponse, nil
}

func (api *API) ListPerformeceYear(ctx context.Context, obj structs.TotaisAnoPerformanceRequest) (structs.TotaisAnoPerformaceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_performance_ano")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.TotaisAnoPerformaceResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return structs.TotaisAnoPerformaceResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.TotaisAnoPerformaceResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return structs.TotaisAnoPerformaceResponse{}, err
	}

	if res == nil {
		return structs.TotaisAnoPerformaceResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.TotaisAnoPerformaceResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return structs.TotaisAnoPerformaceResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse structs.TotaisAnoPerformaceResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.TotaisAnoPerformaceResponse{}, err
	}

	return objResponse, nil
}

func (api *API) ListLifeCostYear(ctx context.Context, obj structs.TotaisAnoCustoDeVidaRequest) (structs.TotaisAnoCustoDeVidaResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_custo_vida_ano")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, err
	}

	if res == nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.TotaisAnoCustoDeVidaResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return structs.TotaisAnoCustoDeVidaResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse structs.TotaisAnoCustoDeVidaResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.TotaisAnoCustoDeVidaResponse{}, err
	}

	return objResponse, nil
}

func (api *API) ListSavingYear(ctx context.Context, obj structs.TotaisAnoEconomizadoRequest) (structs.TotaisAnoEconomizadoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_economizado_ano")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.TotaisAnoEconomizadoResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return structs.TotaisAnoEconomizadoResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.TotaisAnoEconomizadoResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return structs.TotaisAnoEconomizadoResponse{}, err
	}

	if res == nil {
		return structs.TotaisAnoEconomizadoResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.TotaisAnoEconomizadoResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return structs.TotaisAnoEconomizadoResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse structs.TotaisAnoEconomizadoResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.TotaisAnoEconomizadoResponse{}, err
	}

	return objResponse, nil
}

func (api *API) ListDailyAverageYear(ctx context.Context, obj structs.TotaisDiarioMedioAnoRequest) (structs.TotaisDiarioMedioAnoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_diario_medio_ano")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.TotaisDiarioMedioAnoResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return structs.TotaisDiarioMedioAnoResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return structs.TotaisDiarioMedioAnoResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return structs.TotaisDiarioMedioAnoResponse{}, err
	}

	if res == nil {
		return structs.TotaisDiarioMedioAnoResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.TotaisDiarioMedioAnoResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return structs.TotaisDiarioMedioAnoResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse structs.TotaisDiarioMedioAnoResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.TotaisDiarioMedioAnoResponse{}, err
	}

	return objResponse, nil
}
