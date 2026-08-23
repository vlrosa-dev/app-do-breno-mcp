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

func (api *API) ListTags(ctx context.Context, obj structs.TagListRequest) ([]structs.TagListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.BaseUrl, "/rest/v1/rpc/get_totais_tag_mes_v3")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return []structs.TagListResponse{}, err
	}

	bodyStruct, err := json.Marshal(obj)
	if err != nil {
		return []structs.TagListResponse{}, fmt.Errorf("serializando payload: %w", err)
	}
	body := bytes.NewReader(bodyStruct)

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return []structs.TagListResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.ApiToken)
	req.Header.Add("ApiKey", api.ApiKey)

	res, err := api.Client.Do(req)
	if err != nil {
		return []structs.TagListResponse{}, err
	}

	if res == nil {
		return []structs.TagListResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []structs.TagListResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		return []structs.TagListResponse{}, fmt.Errorf(
			"error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData,
		)
	}

	var objResponse []structs.TagListResponse
	err = json.Unmarshal(responseData, &objResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return []structs.TagListResponse{}, err
	}

	return objResponse, nil
}
