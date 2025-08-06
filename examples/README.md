# Ollama Go Examples

This directory contains examples demonstrating how to use the Ollama Go client library.

## Prerequisites

1. Install and run Ollama: [https://ollama.ai/](https://ollama.ai/)
2. Pull a model (e.g., `ollama pull llama2` or `ollama pull gemma3`)

## Examples

### Basic Generation
- **File**: `generate/main.go`
- **Description**: Simple text generation
- **Usage**: `go run examples/generate/main.go`

### Chat
- **File**: `chat/main.go` 
- **Description**: Chat conversation with a model
- **Usage**: `go run examples/chat/main.go`

### Streaming Generation
- **File**: `generate-stream/main.go`
- **Description**: Streaming text generation with real-time output
- **Usage**: `go run examples/generate-stream/main.go`

### Streaming Chat
- **File**: `chat-stream/main.go`
- **Description**: Streaming chat conversation
- **Usage**: `go run examples/chat-stream/main.go`

### Embeddings
- **File**: `embeddings/main.go`
- **Description**: Generate text embeddings
- **Usage**: `go run examples/embeddings/main.go`

### List Models
- **File**: `list/main.go`
- **Description**: List all available models
- **Usage**: `go run examples/list/main.go`

## Running Examples

1. Make sure Ollama is running:
   ```bash
   ollama serve
   ```

2. Pull a model (if you haven't already):
   ```bash
   ollama pull llama2
   # or
   ollama pull gemma3
   ```

3. Run any example:
   ```bash
   go run examples/generate/main.go
   ```

## Modifying Examples

You can modify the examples to use different models by changing the model name in the code:

```go
// Change this line
response, err := ollama.Generate(ctx, "gemma3", "Why is the sky blue?")

// To use a different model
response, err := ollama.Generate(ctx, "llama2", "Why is the sky blue?")
```