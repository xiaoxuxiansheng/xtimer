package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"
)

// 单次读取限制 4M
const (
	defaultReadLimitBytes                = 4 * 1024 * 1024
	defaultTimeoutDuration time.Duration = 5 * time.Second
)

type JSONClient struct {
	timeoutDuration time.Duration
	readLimitBytes  int64
}

func NewJSONClient(opts ...Option) *JSONClient {
	j := JSONClient{}
	for _, opt := range opts {
		opt(&j)
	}

	repair(&j)
	return &j
}

func (j *JSONClient) Get(ctx context.Context, url string, header http.Header, params map[string]string, resp interface{}) error {
	return j.Do(ctx, http.MethodGet, getCompleteURL(url, params), header, nil, resp)
}

func (j *JSONClient) Post(ctx context.Context, url string, header http.Header, req, resp interface{}) error {
	return j.Do(ctx, http.MethodPost, url, header, req, resp)
}

func (j *JSONClient) Patch(ctx context.Context, url string, header http.Header, req, resp interface{}) error {
	return j.Do(ctx, http.MethodPatch, url, header, req, resp)
}

func (j *JSONClient) Delete(ctx context.Context, url string, header http.Header, req, resp interface{}) error {
	return j.Do(ctx, http.MethodDelete, url, header, req, resp)
}

func (j *JSONClient) Do(ctx context.Context, method string, url string, header http.Header, req, resp interface{}) error {
	tCtx, cancel := context.WithTimeout(ctx, j.timeoutDuration)
	defer cancel()

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(tCtx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	for k, vs := range header {
		for _, v := range vs {
			request.Header.Add(k, v)
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(io.LimitReader(response.Body, j.readLimitBytes))
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, resp)
}

func getCompleteURL(originURL string, params map[string]string) string {
	values := neturl.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	queriesStr, _ := neturl.QueryUnescape(values.Encode())
	if len(queriesStr) == 0 {
		return originURL
	}
	return fmt.Sprintf("%s?%s", originURL, queriesStr)
}
