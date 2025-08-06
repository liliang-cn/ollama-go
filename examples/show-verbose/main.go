package main

import (
	"context"
	"fmt"
	"log"

	ollama "github.com/liliang-cn/ollama-go"
)

func main() {
	client, err := ollama.NewClient()
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	ctx := context.Background()
	
	// Show model information with verbose output
	showReq := &ollama.ShowRequest{
		Model:   "llama3.2:1b",
		Verbose: ollama.BoolPtr(true), // Request verbose output
	}
	
	modelInfo, err := client.Show(ctx, showReq)
	if err != nil {
		log.Fatal("Failed to show model:", err)
	}

	fmt.Printf("Model: %s\n", showReq.Model)
	fmt.Printf("Format: %s\n", modelInfo.Details.Format)
	fmt.Printf("Family: %s\n", modelInfo.Details.Family)
	fmt.Printf("Parameter Size: %s\n", modelInfo.Details.ParameterSize)
	if len(modelInfo.Capabilities) > 0 {
		fmt.Printf("Capabilities: %v\n", modelInfo.Capabilities)
	}
	if modelInfo.Template != "" {
		fmt.Printf("Template: %s\n", modelInfo.Template[:100]) // First 100 chars
	}
}