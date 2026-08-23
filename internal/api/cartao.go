package api

import (
	"app-do-breno-agent/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (api *API) ListCards(ctx context.Context, obj structs.ListCartaoRequest) ([]structs.ListCartaoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/cartao")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return []structs.ListCartaoResponse{}, err
	}

	params := url.Values{}
	params.Add("select", obj.Select)
	params.Add("user_id", fmt.Sprintf("eq.%s", obj.UserId))
	params.Add("order", "nome.asc")
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return []structs.ListCartaoResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := api.Client.Do(req)
	if err != nil {
		return []structs.ListCartaoResponse{}, err
	}

	if res == nil {
		return []structs.ListCartaoResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []structs.ListCartaoResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return []structs.ListCartaoResponse{}, msgError
	}

	var despesasResponse []structs.ListCartaoResponse
	err = json.Unmarshal(responseData, &despesasResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return []structs.ListCartaoResponse{}, err
	}

	return despesasResponse, nil
}
