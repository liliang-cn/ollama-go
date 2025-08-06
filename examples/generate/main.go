package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Simple generation
	response, err := ollama.Generate(ctx, "gemma3", "Why is the sky blue?")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Response)
}
