package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// List models
	response, err := ollama.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available models:")
	for _, model := range response.Models {
		fmt.Printf("- %s (size: %d bytes)\n", model.Model, model.Size)
	}
}