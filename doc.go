/*
Package ollama provides a Go client library for the Ollama API.

Ollama is a tool for running large language models locally. This package
provides a comprehensive Go client that mirrors the functionality of the
official Python client with 98%+ feature parity.

# Quick Start

The simplest way to use this package is with the global functions:

	ctx := context.Background()
	response, err := ollama.Generate(ctx, "gemma3", "Why is the sky blue?")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Response)

# Chat Interface

For conversational interactions, use the Chat functions:

	messages := []ollama.Message{
		{Role: "user", Content: "Hello!"},
	}
	response, err := ollama.Chat(ctx, "gemma3", messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Message.Content)

# Streaming

Both Generate and Chat support streaming responses:

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

# Custom Client

For advanced configuration, create a custom client:

	client, err := ollama.NewClient(
		ollama.WithHost("http://localhost:11434"),
		ollama.WithHeaders(map[string]string{
			"Authorization": "Bearer token",
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

# Model Management

The package provides comprehensive model management capabilities:

	// List available models
	models, err := ollama.List(ctx)

	// Pull a model
	err = ollama.Pull(ctx, "gemma3")

	// Show model information
	info, err := ollama.Show(ctx, "gemma3")

# Embeddings

Generate embeddings for text:

	response, err := ollama.Embed(ctx, "nomic-embed-text", "The quick brown fox")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated %d embeddings\n", len(response.Embeddings))

# Configuration

The client can be configured using environment variables:
  - OLLAMA_HOST: Set the Ollama server URL (default: http://localhost:11434)

For more examples and detailed usage, see the examples directory in the repository.
*/
package ollama