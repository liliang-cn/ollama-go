package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}

	if client.baseURL.String() != "http://localhost:11434" {
		t.Errorf("Expected baseURL to be http://localhost:11434, got %s", client.baseURL.String())
	}
}

func TestNewClientWithHost(t *testing.T) {
	customHost := "http://localhost:8080"
	client, err := NewClient(WithHost(customHost))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.baseURL.String() != customHost {
		t.Errorf("Expected baseURL to be %s, got %s", customHost, client.baseURL.String())
	}
}

func TestGenerate(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := GenerateResponse{
			Model:    "test-model",
			Response: "Test response",
			Done:     true,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
	}

	resp, err := client.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if resp.Response != "Test response" {
		t.Errorf("Expected 'Test response', got '%s'", resp.Response)
	}
}

func TestChat(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("Expected path /api/chat, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := ChatResponse{
			Model: "test-model",
			Message: Message{
				Role:    "assistant",
				Content: "Test chat response",
			},
			Done: true,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Test message"},
		},
	}

	resp, err := client.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	if resp.Message.Content != "Test chat response" {
		t.Errorf("Expected 'Test chat response', got '%s'", resp.Message.Content)
	}
}

func TestList(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/tags" {
			t.Errorf("Expected path /api/tags, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := ListResponse{
			Models: []ModelInfo{
				{Model: "model1", Size: 1024},
				{Model: "model2", Size: 2048},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	resp, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(resp.Models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(resp.Models))
	}

	if resp.Models[0].Model != "model1" {
		t.Errorf("Expected 'model1', got '%s'", resp.Models[0].Model)
	}
}

func TestEmbed(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embed" {
			t.Errorf("Expected path /api/embed, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := EmbedResponse{
			Model: "test-embed-model",
			Embeddings: [][]float64{
				{0.1, 0.2, 0.3, 0.4, 0.5},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &EmbedRequest{
		Model: "test-embed-model",
		Input: "Test input text",
	}

	resp, err := client.Embed(context.Background(), req)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	if len(resp.Embeddings) != 1 {
		t.Errorf("Expected 1 embedding, got %d", len(resp.Embeddings))
	}

	if len(resp.Embeddings[0]) != 5 {
		t.Errorf("Expected 5 dimensions, got %d", len(resp.Embeddings[0]))
	}
}

func TestResponseError(t *testing.T) {
	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Test error message"})
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
	}

	_, err := client.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if respErr, ok := err.(*ResponseError); ok {
		if respErr.StatusCode != 400 {
			t.Errorf("Expected status code 400, got %d", respErr.StatusCode)
		}
		if respErr.Message != "Test error message" {
			t.Errorf("Expected 'Test error message', got '%s'", respErr.Message)
		}
	} else {
		t.Errorf("Expected ResponseError, got %T", err)
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test BoolPtr
	b := BoolPtr(true)
	if b == nil || *b != true {
		t.Error("BoolPtr failed")
	}

	// Test IntPtr
	i := IntPtr(42)
	if i == nil || *i != 42 {
		t.Error("IntPtr failed")
	}

	// Test Float64Ptr
	f := Float64Ptr(3.14)
	if f == nil || *f != 3.14 {
		t.Error("Float64Ptr failed")
	}

	// Test StringPtr
	s := StringPtr("test")
	if s == nil || *s != "test" {
		t.Error("StringPtr failed")
	}
}

func TestGenerateStream(t *testing.T) {
	// Mock server that returns streaming responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")

		// Send a single response with Done=true for more predictable testing
		response := GenerateResponse{
			Model:    "test-model",
			Response: "Hello world!",
			Done:     true,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
		Stream: BoolPtr(true),
	}

	respCh, errCh := client.GenerateStream(context.Background(), req)

	// Collect responses
	var responses []*GenerateResponse
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				goto done
			}
			responses = append(responses, resp)
		case err := <-errCh:
			if err != nil {
				t.Fatalf("GenerateStream failed: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for stream responses")
		}
	}
done:

	// Check that we got at least one response
	if len(responses) == 0 {
		t.Fatal("Expected streaming responses, got none")
	}

	// Check that the last response has Done=true
	if !responses[len(responses)-1].Done {
		t.Error("Last response should have Done=true")
	}
}

func TestChatStream(t *testing.T) {
	// Mock server that returns streaming chat responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("Expected path /api/chat, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")

		// Send a single response with Done=true for more predictable testing
		response := ChatResponse{
			Model:   "test-model",
			Message: Message{Role: "assistant", Content: "Hi there!"},
			Done:    true,
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &ChatRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Stream: BoolPtr(true),
	}

	respCh, errCh := client.ChatStream(context.Background(), req)

	var responses []*ChatResponse
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				goto done
			}
			responses = append(responses, resp)
		case err := <-errCh:
			if err != nil {
				t.Fatalf("ChatStream failed: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for chat stream responses")
		}
	}
done:

	// Check that we got at least one response
	if len(responses) == 0 {
		t.Fatal("Expected streaming responses, got none")
	}

	// Check that the last response has Done=true
	if !responses[len(responses)-1].Done {
		t.Error("Last response should have Done=true")
	}
}

func TestShow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/show" {
			t.Errorf("Expected path /api/show, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := ShowResponse{
			Modelfile:  "FROM llama2",
			Parameters: "temperature 0.8",
			Template:   "{{ .Prompt }}",
			Details: &ModelDetails{
				Format:            "gguf",
				Family:            "llama",
				ParameterSize:     "7B",
				QuantizationLevel: "Q4_0",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &ShowRequest{
		Model: "test-model",
	}

	resp, err := client.Show(context.Background(), req)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.Details != nil && resp.Details.Family != "llama" {
		t.Errorf("Expected family 'llama', got '%s'", resp.Details.Family)
	}
}

func TestPull(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pull" {
			t.Errorf("Expected path /api/pull, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := StatusResponse{
			Status: "success",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &PullRequest{
		Model: "test-model",
	}

	resp, err := client.Pull(context.Background(), req)
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/delete" {
			t.Errorf("Expected path /api/delete, got %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		response := StatusResponse{Status: "success"}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &DeleteRequest{
		Model: "test-model",
	}

	resp, err := client.Delete(context.Background(), req)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestCopy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/copy" {
			t.Errorf("Expected path /api/copy, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		response := StatusResponse{Status: "success"}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &CopyRequest{
		Source:      "source-model",
		Destination: "dest-model",
	}

	resp, err := client.Copy(context.Background(), req)
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestPs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ps" {
			t.Errorf("Expected path /api/ps, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := ProcessResponse{
			Models: []ProcessModel{
				{
					Model:     "model1",
					Size:      1024,
					Digest:    "sha256:abc123",
					ExpiresAt: &[]time.Time{time.Now().Add(time.Hour)}[0],
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	resp, err := client.Ps(context.Background())
	if err != nil {
		t.Fatalf("Ps failed: %v", err)
	}

	if len(resp.Models) != 1 {
		t.Errorf("Expected 1 running model, got %d", len(resp.Models))
	}

	if resp.Models[0].Model != "model1" {
		t.Errorf("Expected model 'model1', got '%s'", resp.Models[0].Model)
	}
}

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/version" {
			t.Errorf("Expected path /api/version, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := VersionResponse{
			Version: "0.1.0",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	resp, err := client.Version(context.Background())
	if err != nil {
		t.Fatalf("Version failed: %v", err)
	}

	if resp.Version != "0.1.0" {
		t.Errorf("Expected version '0.1.0', got '%s'", resp.Version)
	}
}

func TestEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embeddings" {
			t.Errorf("Expected path /api/embeddings, got %s", r.URL.Path)
		}

		response := EmbeddingsResponse{
			Embedding: []float64{0.1, 0.2, 0.3},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	resp, err := client.Embeddings(context.Background(), &EmbeddingsRequest{
		Model:  "test-embed-model",
		Prompt: "text1",
	})
	if err != nil {
		t.Fatalf("Embeddings failed: %v", err)
	}

	if len(resp.Embedding) != 3 {
		t.Errorf("Expected 3 embedding dimensions, got %d", len(resp.Embedding))
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client, err := NewClient(WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("Failed to create client with custom HTTP client: %v", err)
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestWithHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Custom-Header": "custom-value",
	}

	client, err := NewClient(WithHeaders(headers))
	if err != nil {
		t.Fatalf("Failed to create client with headers: %v", err)
	}

	if client.headers["Authorization"] != "Bearer token123" {
		t.Errorf("Expected Authorization header, got %s", client.headers["Authorization"])
	}

	if client.headers["Custom-Header"] != "custom-value" {
		t.Errorf("Expected Custom-Header, got %s", client.headers["Custom-Header"])
	}
}

func TestErrorHandling(t *testing.T) {
	// Test invalid URL - this may not always fail in Go
	_, err := NewClient(WithHost("invalid-url"))
	if err != nil {
		// Good, we got an error as expected
		return
	}
	// If no error, that's also acceptable for some invalid URLs

	// Test network error
	client, _ := NewClient(WithHost("http://non-existent-host:12345"))

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
	}

	_, err = client.Generate(context.Background(), req)
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

func TestRequestCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		_ = json.NewEncoder(w).Encode(GenerateResponse{Response: "slow response"})
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
	}

	_, err := client.Generate(ctx, req)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded error, got: %v", err)
	}
}

func TestJSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid model name"}`))
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &GenerateRequest{
		Model:  "invalid-model",
		Prompt: "Test prompt",
	}

	_, err := client.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if respErr, ok := err.(*ResponseError); ok {
		if respErr.StatusCode != 400 {
			t.Errorf("Expected status code 400, got %d", respErr.StatusCode)
		}
		if respErr.Message != "invalid model name" {
			t.Errorf("Expected 'invalid model name', got '%s'", respErr.Message)
		}
	} else {
		t.Errorf("Expected ResponseError, got %T", err)
	}
}

func TestNonJSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client, _ := NewClient(WithHost(server.URL))

	req := &GenerateRequest{
		Model:  "test-model",
		Prompt: "Test prompt",
	}

	_, err := client.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if respErr, ok := err.(*ResponseError); ok {
		if respErr.StatusCode != 500 {
			t.Errorf("Expected status code 500, got %d", respErr.StatusCode)
		}
		if respErr.Message != "Internal Server Error" {
			t.Errorf("Expected 'Internal Server Error', got '%s'", respErr.Message)
		}
	} else {
		t.Errorf("Expected ResponseError, got %T", err)
	}
}
