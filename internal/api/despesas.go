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

func (api *API) HandleCreateTransactionHTTP(ctx context.Context, obj *structs.TransactionRequest) (string, error) {
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

func (api *API) HandleDeleteTransactionHTTP(ctx context.Context, obj *structs.TransactionDeleteRequest) error {
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

func (api *API) HandleCreateTransactionInSerieHTTP(ctx context.Context, obj structs.TransactionInSerieRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/criar_movimentacoes_repetidas_v2")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if res.StatusCode > http.StatusNoContent {
		return fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	return nil
}

func (api *API) HandleDeleteTransactionInSerieHTTP(ctx context.Context, obj structs.TransactionInSerieDeleteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/serie_delete_from")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if res.StatusCode > http.StatusNoContent {
		return fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	return nil
}

func (api *API) HandleDeleteAllTransactionInSerieHTTP(ctx context.Context, obj structs.TransactionInSerieDeleteAllRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/serie_delete_all")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

func (api *API) HandleFetchTransactionsHTTP(ctx context.Context, obj structs.TransactionListRequest) ([]structs.TransactionListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/view_movimentacao_all")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}

	params := url.Values{}
	params.Add("select", obj.Select)
	params.Add("user_id", fmt.Sprintf("eq.%s", obj.UserId))
	params.Add("tipo", fmt.Sprintf("eq.%s", obj.Tipo))
	params.Add("data", fmt.Sprintf("gte.%s", obj.DataInicial))
	params.Add("data", fmt.Sprintf("lt.%s", obj.DataFinal))
	params.Add("order", "data.desc, created_at.desc")
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.Client.Do(req)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}

	if res == nil {
		return []structs.TransactionListResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return []structs.TransactionListResponse{}, msgError
	}

	var despesasResponse []structs.TransactionListResponse
	err = json.Unmarshal(responseData, &despesasResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return []structs.TransactionListResponse{}, err
	}

	return despesasResponse, nil
}

func (api *API) HandleGetTransactionByIdHTTP(ctx context.Context, obj structs.TransactionGetRequest) ([]structs.TransactionListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/view_movimentacao_all")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}

	params := url.Values{}
	params.Add("select", obj.Select)
	params.Add("limit", fmt.Sprintf("%d", obj.Limit))
	params.Add("movimentacao_id", fmt.Sprintf("eq.%s", obj.MovimentacaoId))
	params.Add("user_id", fmt.Sprintf("eq.%s", api.UserID))
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.Client.Do(req)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}

	if res == nil {
		return []structs.TransactionListResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []structs.TransactionListResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return []structs.TransactionListResponse{}, msgError
	}

	var despesasResponse []structs.TransactionListResponse
	err = json.Unmarshal(responseData, &despesasResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return []structs.TransactionListResponse{}, err
	}

	return despesasResponse, nil
}

func (api *API) HandlePutTransactionHTTP(ctx context.Context, obj structs.TransactionUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/movimentacao_update_with_tags_v2")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if res.StatusCode > http.StatusNoContent {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return msgError
	}

	return nil
}
