package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func MakeHTTPRequest[T any](
	c *http.Client,
	ctx context.Context,
	fullUrl string,
	httpMethod string,
	headers map[string]string,
	queryParameters url.Values,
	body io.Reader,
	responseType T,
) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client := c
	u, err := url.Parse(fullUrl)
	if err != nil {
		return responseType, err
	}

	// set the query to the encoded parameters
	u.RawQuery = queryParameters.Encode()

	// regardless of GET or POST, we can safely add the body
	req, err := http.NewRequestWithContext(ctx, httpMethod, u.String(), body)
	if err != nil {
		return responseType, err
	}

	// for each header passed, add the header value to the request
	for k, v := range headers {
		req.Header.Set(k, v)
		// log.Printf("Headers:%s:%v", k, v)
	}

	// optional: log the request for easier stack tracing
	log.Printf("%s %s\n", httpMethod, req.URL.String())

	// finally, do the request
	res, err := client.Do(req)
	if err != nil {
		return responseType, err
	}

	if res == nil {
		return responseType, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return responseType, err
	}

	defer res.Body.Close()

	if res.StatusCode > http.StatusNoContent {
		msgError := fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
		log.Println(msgError)
		return responseType, msgError
	}

	var responseObject T
	// some endpoints (e.g. rpc calls) return a raw string body (like a UUID)
	// instead of JSON, so json.Unmarshal would fail on it — assign it directly.
	if s, ok := any(&responseObject).(*string); ok {
		*s = string(responseData)
		return responseObject, nil
	}

	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Printf("error unmarshaling response: %+v", err)
		return responseType, err
	}

	return responseObject, nil
}
