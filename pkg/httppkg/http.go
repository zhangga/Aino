package httppkg

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	json "github.com/bytedance/sonic"
)

type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Body    any // 自动 JSON marshal
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type DefaultClient struct {
	httpClient *http.Client
	baseHeader http.Header
}

// NewDefaultClient UnThreadsafe
func NewDefaultClient(timeout time.Duration) *DefaultClient {
	return &DefaultClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseHeader: make(http.Header),
	}
}

func (c *DefaultClient) SetHeader(key, value string) {
	c.baseHeader.Set(key, value)
}

func (c *DefaultClient) Do(ctx context.Context, r *Request) (*Response, error) {
	if r.Method == "" {
		return nil, errors.New("http method is required")
	}
	if r.URL == "" {
		return nil, errors.New("url is required")
	}

	var bodyReader io.Reader
	if r.Body != nil {
		data, err := json.Marshal(r.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// base headers
	for k, v := range c.baseHeader {
		req.Header[k] = v
	}

	// request headers
	for k, v := range r.Headers {
		req.Header[k] = v
	}

	// default content-type
	if r.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}, nil
}
