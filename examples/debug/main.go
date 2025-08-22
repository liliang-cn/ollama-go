package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// First, list available models
	fmt.Println("Available models:")
	models, err := ollama.List(ctx)
	if err != nil {
		log.Fatal("Failed to list models:", err)
	}

	if len(models.Models) == 0 {
		fmt.Println("No models found. Please pull a model first:")
		fmt.Println("  ollama pull llama3.2")
		fmt.Println("  ollama pull qwen2.5")
		return
	}

	for i, model := range models.Models {
		fmt.Printf("%d. %s\n", i+1, model.Model)
	}

	// Use the first available model for generation
	modelName := models.Models[0].Model
	fmt.Printf("\nUsing model: %s\n", modelName)

	// Simple generation with detailed output
	fmt.Println("\nGenerating response...")
	response, err := ollama.Generate(ctx, modelName, "Why is the sky blue? Please give a short answer.")
	if err != nil {
		log.Fatal("Generation failed:", err)
	}

	fmt.Printf("\nResponse: %s\n", response.Response)
	if response.Response == "" {
		fmt.Println("Warning: Empty response received")
		fmt.Printf("Model: %s\n", response.Model)
		fmt.Printf("Done: %t\n", response.Done)
		fmt.Printf("Total Duration: %d ns\n", response.TotalDuration)
	}
}
