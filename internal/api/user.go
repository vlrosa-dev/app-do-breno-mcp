package api

import (
	"app-do-breno-agent/internal/helpers"
	"app-do-breno-agent/internal/structs"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (api *API) GetUserInfo(ctx context.Context) (structs.UserInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", api.baseURL, "/auth/v1/user")

	headers := map[string]string{
		"Authorization": "Bearer " + api.token,
		"ApiKey":        api.apiKey,
	}

	params := url.Values{}

	var userInfoResponse structs.UserInfoResponse
	resp, err := helpers.MakeHTTPRequest(
		api.client,
		ctx,
		fullUrl,
		http.MethodGet,
		headers,
		params,
		nil,
		userInfoResponse,
	)
	if err != nil {
		return structs.UserInfoResponse{}, fmt.Errorf("request error, msg: %v", err.Error())
	}

	return resp, err
}
