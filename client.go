package beacon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is a beacon node API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new beacon node API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{},
	}
}

// APIError represents an error response from the beacon node API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("beacon API error (code %d): %s", e.Code, e.Message)
}

// doRequest performs an HTTP request and returns the raw response body
// Each endpoint should define its own response structure and unmarshal accordingly
func (c *Client) doRequest(ctx context.Context, method, endpoint string, query url.Values) ([]byte, error) {
	fullURL := c.baseURL + endpoint
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
		}
		return nil, &apiErr
	}

	return body, nil
}
