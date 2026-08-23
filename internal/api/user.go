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

func GetUser(ctx context.Context, client *http.Client, baseUrl, apiToken, apiKey string) (structs.UserInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fullUrl := fmt.Sprintf("%s%s", baseUrl, "/auth/v1/user")
	u, err := url.Parse(fullUrl)
	if err != nil {
		return structs.UserInfoResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return structs.UserInfoResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("ApiKey", apiKey)

	res, err := client.Do(req)
	if err != nil {
		return structs.UserInfoResponse{}, err
	}

	if res == nil {
		return structs.UserInfoResponse{}, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return structs.UserInfoResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return structs.UserInfoResponse{}, msgError
	}

	var userInfoResponse structs.UserInfoResponse
	err = json.Unmarshal(responseData, &userInfoResponse)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return structs.UserInfoResponse{}, err
	}

	return userInfoResponse, nil
}
