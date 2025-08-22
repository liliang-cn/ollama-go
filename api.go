package ollama

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Generate generates a response from a prompt
func (c *Client) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// Ensure stream is set to false for non-streaming requests
	generateReq := *req
	generateReq.Stream = BoolPtr(false)

	resp, err := c.doRequest(ctx, "POST", "/api/generate", &generateReq)
	if err != nil {
		return nil, err
	}

	var result GenerateResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse generate response: %w", err)
	}

	// Process <think> tags in response if present
	if cleanResponse, thinking := extractThinkingContent(result.Response); thinking != "" {
		result.Response = cleanResponse
		// Only set thinking if it wasn't already set by the server
		if result.Thinking == "" {
			result.Thinking = thinking
		}
	}

	return &result, nil
}

// GenerateStream generates a streaming response from a prompt
func (c *Client) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan *GenerateResponse, <-chan error) {
	responseChan := make(chan *GenerateResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		streamReq := *req
		streamReq.Stream = BoolPtr(true)

		resp, err := c.doRequest(ctx, "POST", "/api/generate", &streamReq)
		if err != nil {
			errorChan <- err
			return
		}

		dataChan, errChan := c.parseStreamResponse(resp)
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}
				var genResp GenerateResponse
				if err := json.Unmarshal(data, &genResp); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}

				// Process <think> tags in streaming response if present
				if cleanResponse, thinking := extractThinkingContent(genResp.Response); thinking != "" {
					genResp.Response = cleanResponse
					// Only set thinking if it wasn't already set by the server
					if genResp.Thinking == "" {
						genResp.Thinking = thinking
					}
				}

				responseChan <- &genResp
				if genResp.Done {
					return
				}
			case err := <-errChan:
				if err != nil {
					errorChan <- err
				}
				return
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// Chat sends a chat request and returns the response
func (c *Client) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Ensure stream is set to false for non-streaming requests
	chatReq := *req
	chatReq.Stream = BoolPtr(false)

	resp, err := c.doRequest(ctx, "POST", "/api/chat", &chatReq)
	if err != nil {
		return nil, err
	}

	var result ChatResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse chat response: %w", err)
	}

	// Process <think> tags in message content if present
	if cleanContent, thinking := extractThinkingContent(result.Message.Content); thinking != "" {
		result.Message.Content = cleanContent
		// Only set thinking if it wasn't already set by the server
		if result.Message.Thinking == "" {
			result.Message.Thinking = thinking
		}
	}

	return &result, nil
}

// ChatStream sends a chat request and returns a streaming response
func (c *Client) ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatResponse, <-chan error) {
	responseChan := make(chan *ChatResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		streamReq := *req
		streamReq.Stream = BoolPtr(true)

		resp, err := c.doRequest(ctx, "POST", "/api/chat", &streamReq)
		if err != nil {
			errorChan <- err
			return
		}

		dataChan, errChan := c.parseStreamResponse(resp)
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}
				var chatResp ChatResponse
				if err := json.Unmarshal(data, &chatResp); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}

				// Process <think> tags in streaming message content if present
				if cleanContent, thinking := extractThinkingContent(chatResp.Message.Content); thinking != "" {
					chatResp.Message.Content = cleanContent
					// Only set thinking if it wasn't already set by the server
					if chatResp.Message.Thinking == "" {
						chatResp.Message.Thinking = thinking
					}
				}

				responseChan <- &chatResp
				if chatResp.Done {
					return
				}
			case err := <-errChan:
				if err != nil {
					errorChan <- err
				}
				return
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// Embed creates embeddings for the given input
func (c *Client) Embed(ctx context.Context, req *EmbedRequest) (*EmbedResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/embed", req)
	if err != nil {
		return nil, err
	}

	var result EmbedResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse embed response: %w", err)
	}

	return &result, nil
}

// Embeddings creates embeddings for the given prompt (legacy API)
func (c *Client) Embeddings(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/embeddings", req)
	if err != nil {
		return nil, err
	}

	var result EmbeddingsResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse embeddings response: %w", err)
	}

	return &result, nil
}

// List lists available models
func (c *Client) List(ctx context.Context) (*ListResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/tags", nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse list response: %w", err)
	}

	return &result, nil
}

// Show returns information about a model
func (c *Client) Show(ctx context.Context, req *ShowRequest) (*ShowResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/show", req)
	if err != nil {
		return nil, err
	}

	var result ShowResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse show response: %w", err)
	}

	return &result, nil
}

// Pull downloads a model
func (c *Client) Pull(ctx context.Context, req *PullRequest) (*StatusResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/pull", req)
	if err != nil {
		return nil, err
	}

	var result StatusResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse pull response: %w", err)
	}

	return &result, nil
}

// PullStream downloads a model with progress updates
func (c *Client) PullStream(ctx context.Context, req *PullRequest) (<-chan *ProgressResponse, <-chan error) {
	responseChan := make(chan *ProgressResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		streamReq := *req
		streamReq.Stream = BoolPtr(true)

		resp, err := c.doRequest(ctx, "POST", "/api/pull", &streamReq)
		if err != nil {
			errorChan <- err
			return
		}

		dataChan, errChan := c.parseStreamResponse(resp)
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}
				var progResp ProgressResponse
				if err := json.Unmarshal(data, &progResp); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}
				responseChan <- &progResp
			case err := <-errChan:
				if err != nil {
					errorChan <- err
				}
				return
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// Push uploads a model
func (c *Client) Push(ctx context.Context, req *PushRequest) (*StatusResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/push", req)
	if err != nil {
		return nil, err
	}

	var result StatusResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse push response: %w", err)
	}

	return &result, nil
}

// PushStream uploads a model with progress updates
func (c *Client) PushStream(ctx context.Context, req *PushRequest) (<-chan *ProgressResponse, <-chan error) {
	responseChan := make(chan *ProgressResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		streamReq := *req
		streamReq.Stream = BoolPtr(true)

		resp, err := c.doRequest(ctx, "POST", "/api/push", &streamReq)
		if err != nil {
			errorChan <- err
			return
		}

		dataChan, errChan := c.parseStreamResponse(resp)
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}
				var progResp ProgressResponse
				if err := json.Unmarshal(data, &progResp); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}
				responseChan <- &progResp
			case err := <-errChan:
				if err != nil {
					errorChan <- err
				}
				return
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// Create creates a new model
func (c *Client) Create(ctx context.Context, req *CreateRequest) (*StatusResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/create", req)
	if err != nil {
		return nil, err
	}

	var result StatusResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	return &result, nil
}

// CreateStream creates a new model with progress updates
func (c *Client) CreateStream(ctx context.Context, req *CreateRequest) (<-chan *ProgressResponse, <-chan error) {
	responseChan := make(chan *ProgressResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		streamReq := *req
		streamReq.Stream = BoolPtr(true)

		resp, err := c.doRequest(ctx, "POST", "/api/create", &streamReq)
		if err != nil {
			errorChan <- err
			return
		}

		dataChan, errChan := c.parseStreamResponse(resp)
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					return
				}
				var progResp ProgressResponse
				if err := json.Unmarshal(data, &progResp); err != nil {
					errorChan <- fmt.Errorf("failed to parse streaming response: %w", err)
					return
				}
				responseChan <- &progResp
			case err := <-errChan:
				if err != nil {
					errorChan <- err
				}
				return
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	return responseChan, errorChan
}

// Delete deletes a model
func (c *Client) Delete(ctx context.Context, req *DeleteRequest) (*StatusResponse, error) {
	resp, err := c.doRequest(ctx, "DELETE", "/api/delete", req)
	if err != nil {
		return nil, err
	}

	status := "success"
	if resp.StatusCode != 200 {
		status = "error"
	}

	return &StatusResponse{Status: status}, nil
}

// Copy copies a model
func (c *Client) Copy(ctx context.Context, req *CopyRequest) (*StatusResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/copy", req)
	if err != nil {
		return nil, err
	}

	status := "success"
	if resp.StatusCode != 200 {
		status = "error"
	}

	return &StatusResponse{Status: status}, nil
}

// Ps shows running processes
func (c *Client) Ps(ctx context.Context) (*ProcessResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/ps", nil)
	if err != nil {
		return nil, err
	}

	var result ProcessResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ps response: %w", err)
	}

	return &result, nil
}

// Version gets the Ollama server version
func (c *Client) Version(ctx context.Context) (*VersionResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/version", nil)
	if err != nil {
		return nil, err
	}

	var result VersionResponse
	if err := c.parseJSONResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse version response: %w", err)
	}

	return &result, nil
}

// BoolPtr returns a pointer to a bool value
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns a pointer to an int value
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to a float64 value
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr returns a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// CreateBlob uploads a file and returns its digest
func (c *Client) CreateBlob(ctx context.Context, path string) (string, error) {
	// Calculate SHA256 hash
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}
	digest := fmt.Sprintf("sha256:%x", hash.Sum(nil))

	// Reopen file for upload
	file, err = os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to reopen file: %w", err)
	}
	defer file.Close()

	// Upload the blob
	resp, err := c.doRequestWithBody(ctx, "POST", fmt.Sprintf("/api/blobs/%s", digest), file)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}
	resp.Body.Close()

	return digest, nil
}

// CheckBlob checks if a blob exists on the server
func (c *Client) CheckBlob(ctx context.Context, digest string) (bool, error) {
	resp, err := c.doRequest(ctx, "HEAD", fmt.Sprintf("/api/blobs/%s", digest), nil)
	if err != nil {
		// Check if it's a 404 error (blob doesn't exist)
		if respErr, ok := err.(*ResponseError); ok && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("failed to check blob: %w", err)
	}
	resp.Body.Close()

	// If we get here, the blob exists (200 OK)
	return true, nil
}
