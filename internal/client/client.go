package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BaseURL string

const (
	BaseTransactional BaseURL = "https://send.api.mailtrap.io"
	BaseBulk          BaseURL = "https://bulk.api.mailtrap.io"
	BaseSandbox       BaseURL = "https://sandbox.api.mailtrap.io"
	BaseGeneral       BaseURL = "https://mailtrap.io"
)

type Client struct {
	httpClient    *http.Client
	apiToken      string
	baseOverrides map[BaseURL]string
}

func New(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiToken:   apiToken,
	}
}

// SetBaseURL overrides the base URL for the given key. Useful in tests.
func (c *Client) SetBaseURL(base BaseURL, rawurl string) {
	if c.baseOverrides == nil {
		c.baseOverrides = make(map[BaseURL]string)
	}
	c.baseOverrides[base] = rawurl
}

func (c *Client) resolveBase(base BaseURL) string {
	if c.baseOverrides != nil {
		if override, ok := c.baseOverrides[base]; ok {
			return override
		}
	}
	return string(base)
}

func (c *Client) do(ctx context.Context, base BaseURL, method, path string, query url.Values, body interface{}, result interface{}) error {
	u := c.resolveBase(base) + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.setAuthHeader(req, base)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if json.Unmarshal(respBody, apiErr) != nil {
			apiErr.Message = string(respBody)
		}
		return apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) setAuthHeader(req *http.Request, base BaseURL) {
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
}

func (c *Client) Get(ctx context.Context, base BaseURL, path string, query url.Values, result interface{}) error {
	return c.do(ctx, base, http.MethodGet, path, query, nil, result)
}

func (c *Client) Post(ctx context.Context, base BaseURL, path string, body interface{}, result interface{}) error {
	return c.do(ctx, base, http.MethodPost, path, nil, body, result)
}

func (c *Client) Patch(ctx context.Context, base BaseURL, path string, body interface{}, result interface{}) error {
	return c.do(ctx, base, http.MethodPatch, path, nil, body, result)
}

func (c *Client) Put(ctx context.Context, base BaseURL, path string, body interface{}, result interface{}) error {
	return c.do(ctx, base, http.MethodPut, path, nil, body, result)
}

func (c *Client) Delete(ctx context.Context, base BaseURL, path string, result interface{}) error {
	return c.do(ctx, base, http.MethodDelete, path, nil, nil, result)
}

// GetRaw returns the raw response body bytes (for eml, raw, etc.)
func (c *Client) GetRaw(ctx context.Context, base BaseURL, path string, query url.Values) ([]byte, error) {
	u := c.resolveBase(base) + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setAuthHeader(req, base)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: string(data)}
	}

	return data, nil
}
