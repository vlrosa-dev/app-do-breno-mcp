package api

import (
	"app-do-breno-agent/internal/helpers"
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

func (api *API) CreateExpenseOne(ctx context.Context, obj *structs.DespesasUnicaRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/movimentacao_create_with_tags_v2")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return "", err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res == nil {
		return "", fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return "", msgError
	}

	return string(responseData), nil
}

func (api *API) DeleteExpenseOne(ctx context.Context, obj *structs.DespesasUnicaDeleteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/movimentacao")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("id", fmt.Sprintf("eq.%s", obj.Id))
	params.Add("user_id", fmt.Sprintf("eq.%s", obj.UserId))
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return err
	}

	if res == nil {
		return fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	return nil
}

func (api *API) CreateExpensesRecurrency(ctx context.Context, obj structs.DespesasRecorrenteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/criar_movimentacoes_repetidas_v2")

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + api.ApiToken,
		"ApiKey":        api.ApiKey,
	}

	params := url.Values{}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	_, err = helpers.MakeHTTPRequest(
		api.Client,
		ctx,
		fullUrl,
		http.MethodPost,
		headers,
		params,
		body,
		struct{}{},
	)
	if err != nil {
		return fmt.Errorf("request error, msg: %v", err.Error())
	}

	return nil
}

func (api *API) DeleteExpensesRecurrency(ctx context.Context, obj structs.DespesasRecorrentesDeleteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/serie_delete_from")

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + api.ApiToken,
		"ApiKey":        api.ApiKey,
	}

	params := url.Values{}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	_, err = helpers.MakeHTTPRequest(
		api.Client,
		ctx,
		fullUrl,
		http.MethodPost,
		headers,
		params,
		body,
		struct{}{},
	)
	if err != nil {
		return fmt.Errorf("request error, msg: %v", err.Error())
	}

	return nil
}

func (api *API) DeleteExpensesRecurrencyAll(ctx context.Context, obj structs.DespesasRecorrentesTodosDeleteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/serie_delete_all")

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + api.ApiToken,
		"ApiKey":        api.ApiKey,
	}

	params := url.Values{}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	_, err = helpers.MakeHTTPRequest(
		api.Client,
		ctx,
		fullUrl,
		http.MethodPost,
		headers,
		params,
		body,
		struct{}{},
	)
	if err != nil {
		return fmt.Errorf("request error, msg: %v", err.Error())
	}

	return nil
}

func (api *API) ListExpenses(ctx context.Context, obj structs.DespesasListarRequest) ([]structs.DespesasListarResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/view_movimentacao_all")

	headers := map[string]string{
		"Authorization": "Bearer " + api.ApiToken,
		"ApiKey":        api.ApiKey,
	}

	params := url.Values{}
	params.Add("select", obj.Select)
	params.Add("user_id", fmt.Sprintf("eq.%s", obj.UserId))
	params.Add("tipo", fmt.Sprintf("eq.%s", obj.Tipo))
	params.Add("data", fmt.Sprintf("gte.%s", obj.DataInicial))
	params.Add("data", fmt.Sprintf("lt.%s", obj.DataFinal))
	params.Add("order", "data.desc, created_at.desc")

	var despesas []structs.DespesasListarResponse
	resp, err := helpers.MakeHTTPRequest(
		api.Client,
		ctx,
		fullUrl,
		http.MethodGet,
		headers,
		params,
		nil,
		despesas,
	)
	if err != nil {
		return []structs.DespesasListarResponse{}, fmt.Errorf("request error, msg: %v", err.Error())
	}

	return resp, err
}
