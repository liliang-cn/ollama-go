package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Create embeddings - matching Python example
	response, err := ollama.Embed(ctx, "llama3.2", "Hello, world!")
	if err != nil {
		log.Fatal(err)
	}

	// Print embeddings like Python example
	fmt.Printf("%v\n", response.Embeddings)
}
