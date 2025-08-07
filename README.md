# Ollama Go Client

A Go client library for [Ollama](https://ollama.ai/), based on the official Python client.

> **Note**: This is an unofficial Go client library inspired by the official Python client. It provides the same functionality and API design patterns with **98%+ feature parity**.

## ✨ Features

- **Complete API Support**: All Ollama REST API endpoints
- **Streaming Support**: Real-time streaming for generation and chat
- **Type Safety**: Full Go type definitions with compile-time checking  
- **Flexible Configuration**: Multiple client configuration options
- **Error Handling**: Comprehensive error handling with JSON parsing
- **File Upload**: Blob upload functionality for model creation
- **Advanced Options**: 20+ configuration functions for fine-tuning
- **Context Support**: Full context.Context support for cancellation

## Installation

```bash
go get github.com/liliang-cn/ollama-go
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    // Generate a response
    response, err := ollama.Generate(ctx, "gemma3", "Why is the sky blue?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Response)
}
```

### Chat

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    messages := []ollama.Message{
        {
            Role:    "user",
            Content: "Why is the sky blue?",
        },
    }

    response, err := ollama.Chat(ctx, "gemma3", messages)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Message.Content)
}
```

### Streaming

Both `Generate` and `Chat` support streaming responses:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    // Stream generation
    responseChan, errorChan := ollama.GenerateStream(ctx, "gemma3", "Tell me a story")

    for {
        select {
        case response, ok := <-responseChan:
            if !ok {
                return
            }
            fmt.Print(response.Response)
        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

### Custom Client

You can create a custom client with specific configuration:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    // Create a custom HTTP client
    httpClient := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create client with custom configuration
    client, err := ollama.NewClient(
        ollama.WithHost("http://localhost:11434"),
        ollama.WithHTTPClient(httpClient),
        ollama.WithHeaders(map[string]string{
            "Custom-Header": "custom-value",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    
    req := &ollama.GenerateRequest{
        Model:  "gemma3",
        Prompt: "Hello, world!",
    }

    response, err := client.Generate(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Response)
}
```

### Embeddings

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    // Create embeddings
    response, err := ollama.Embed(ctx, "nomic-embed-text", "The quick brown fox")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated %d embeddings\n", len(response.Embeddings))
}
```

### Model Creation with All Options

```go
package main

import (
    "context"
    
    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    client, _ := ollama.NewClient()
    
    // Create model with complete configuration
    req := &ollama.CreateRequest{
        Model:     "my-custom-model",
        Modelfile: "FROM llama2\nSYSTEM \"You are a helpful assistant.\"",
        Files:     map[string]string{"data.txt": "training data"},
        Adapters:  map[string]string{"lora": "adapter_data"},
        Template:  "{{ .System }}{{ .Prompt }}",
        License:   "MIT",
        System:    "Custom system prompt",
        Parameters: &ollama.Options{
            Temperature: ollama.Float64Ptr(0.7),
        },
        Messages: []ollama.Message{
            {Role: "system", Content: "You are helpful"},
        },
    }
    
    status, err := client.Create(ctx, req)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Model created: %s\n", status.Status)
}
```

### File Upload (Blob)

```go
package main

import (
    "context"
    
    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()
    
    // Upload a file and get its digest
    digest, err := ollama.CreateBlob(ctx, "/path/to/file.bin")
    if err != nil {
        panic(err)
    }
    fmt.Printf("File uploaded with digest: %s\n", digest)
}
```

### Progress Streaming

For operations like pulling models, you can stream progress updates:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/liliang-cn/ollama-go"
)

func main() {
    ctx := context.Background()

    progressChan, errorChan := ollama.PullStream(ctx, "gemma3")

    for {
        select {
        case progress, ok := <-progressChan:
            if !ok {
                fmt.Println("Pull completed!")
                return
            }
            if progress.Total > 0 {
                percentage := float64(progress.Completed) / float64(progress.Total) * 100
                fmt.Printf("Progress: %.1f%% (%s)\n", percentage, progress.Status)
            } else {
                fmt.Printf("Status: %s\n", progress.Status)
            }
        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

## API

### Client Methods

- `Generate(ctx, req)` - Generate a completion
- `GenerateStream(ctx, req)` - Generate a streaming completion
- `Chat(ctx, req)` - Send a chat message
- `ChatStream(ctx, req)` - Send a chat message with streaming response
- `Embed(ctx, req)` - Create embeddings
- `Embeddings(ctx, req)` - Create embeddings (legacy API)
- `List(ctx)` - List available models
- `Show(ctx, req)` - Show model information
- `Pull(ctx, req)` - Download a model
- `PullStream(ctx, req)` - Download a model with progress
- `Push(ctx, req)` - Upload a model
- `PushStream(ctx, req)` - Upload a model with progress
- `Create(ctx, req)` - Create a model from a Modelfile
- `CreateStream(ctx, req)` - Create a model with progress
- `Delete(ctx, req)` - Delete a model
- `Copy(ctx, req)` - Copy a model
- `Ps(ctx)` - List running processes

### Global Functions

For convenience, all client methods are also available as global functions that use a default client instance:

- `ollama.Generate(ctx, model, prompt, options...)`
- `ollama.Chat(ctx, model, messages, options...)`
- `ollama.Embed(ctx, model, input, options...)`
- And so on...

### Configuration Options

The client can be configured using option functions:

- `WithHost(host)` - Set the Ollama server URL
- `WithHTTPClient(client)` - Use a custom HTTP client
- `WithHeaders(headers)` - Add custom headers

### Request Options

Many functions support option functions for common configurations:

- `WithOptions(options)` - Set model options
- `WithSystem(prompt)` - Set system prompt
- `WithFormat(format)` - Set response format
- `WithKeepAlive(duration)` - Set keep alive duration
- `WithImages(images)` - Add images (for multimodal models)
- `WithTools(tools)` - Add tools for function calling
- `WithThinking()` - Enable thinking mode

## Environment Variables

- `OLLAMA_HOST` - Set the Ollama server URL (default: `http://localhost:11434`)

## License

This project is licensed under the MIT License.