package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
		json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Test error message"})
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