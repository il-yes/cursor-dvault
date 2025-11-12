package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type ApiBase struct {
	BaseURL string
	Client  *http.Client
}

func NewApiBase(baseURL string) *ApiBase {
	return &ApiBase{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (a *ApiBase) Get(path string) ([]byte, error) {
	resp, err := a.Client.Get(a.BaseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (a *ApiBase) Post(path string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Post(a.BaseURL+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// You can add PUT, DELETE similarly
