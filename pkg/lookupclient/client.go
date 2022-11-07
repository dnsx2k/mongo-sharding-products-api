package lookupclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClientCtx struct {
	httpClient *http.Client
	baseURL    string
}

func New(baseURL string) *ClientCtx {
	return &ClientCtx{httpClient: http.DefaultClient, baseURL: baseURL}
}

func (cc *ClientCtx) GetLookup(key string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", cc.baseURL, key))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseDto map[string]interface{}
	if err := json.Unmarshal(body, &responseDto); err != nil {
		return "", err
	}

	return responseDto["location"].(string), nil
}
