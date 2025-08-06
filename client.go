package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	defaultHost    = "http://localhost:11434"
	userAgent      = "ollama-go/1.0.0 (liliang-cn)"
	requestTimeout = 120 * time.Second  // Increased for large models
)

// Client represents the Ollama API client
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	headers    map[string]string
}

// ClientOption defines a function type for configuring the client
type ClientOption func(*Client)

// NewClient creates a new Ollama client with optional configuration
func NewClient(options ...ClientOption) (*Client, error) {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = defaultHost
	}

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host URL: %w", err)
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		baseURL: baseURL,
		headers: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   userAgent,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// WithHost sets the Ollama host URL
func WithHost(host string) ClientOption {
	return func(c *Client) {
		if u, err := url.Parse(host); err == nil {
			c.baseURL = u
		}
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithHeaders adds custom headers
func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL.String() + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		// Try to parse as JSON error first, like Python version
		var errorResp map[string]interface{}
		errorMsg := string(bodyBytes)
		if json.Unmarshal(bodyBytes, &errorResp) == nil {
			if errStr, ok := errorResp["error"].(string); ok && errStr != "" {
				errorMsg = errStr
			}
		}
		
		return nil, &ResponseError{
			StatusCode: resp.StatusCode,
			Message:    errorMsg,
		}
	}

	return resp, nil
}

// doRequestWithBody performs an HTTP request with a body reader (for file uploads)
func (c *Client) doRequestWithBody(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.baseURL.String() + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers, but exclude Content-Type for file uploads to let HTTP set it
	for key, value := range c.headers {
		if key != "Content-Type" {
			req.Header.Set(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		// Try to parse as JSON error first, like Python version
		var errorResp map[string]interface{}
		errorMsg := string(bodyBytes)
		if json.Unmarshal(bodyBytes, &errorResp) == nil {
			if errStr, ok := errorResp["error"].(string); ok && errStr != "" {
				errorMsg = errStr
			}
		}
		
		return nil, &ResponseError{
			StatusCode: resp.StatusCode,
			Message:    errorMsg,
		}
	}

	return resp, nil
}

// parseJSONResponse parses a JSON response into the target struct
func (c *Client) parseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// parseStreamResponse handles streaming responses
func (c *Client) parseStreamResponse(resp *http.Response) (<-chan []byte, <-chan error) {
	dataChan := make(chan []byte, 100)
	errChan := make(chan error, 1)

	go func() {
		defer resp.Body.Close()
		defer close(dataChan)
		defer close(errChan)

		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var rawMessage json.RawMessage
			if err := decoder.Decode(&rawMessage); err != nil {
				if err != io.EOF {
					errChan <- fmt.Errorf("failed to decode response: %w", err)
				}
				return
			}
			dataChan <- rawMessage
		}
	}()

	return dataChan, errChan
}