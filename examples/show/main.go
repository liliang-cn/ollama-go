package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Show model information
	response, err := ollama.Show(ctx, "gemma3")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Model Information:")
	fmt.Printf("Modified at:   %s\n", response.ModifiedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Template:      %s\n", response.Template)
	fmt.Printf("Modelfile:     %s\n", response.Modelfile)
	fmt.Printf("License:       %s\n", response.License)

	if response.Details != nil {
		fmt.Printf("Details:\n")
		fmt.Printf("  Format:      %s\n", response.Details.Format)
		fmt.Printf("  Family:      %s\n", response.Details.Family)
		fmt.Printf("  Parameter Size: %s\n", response.Details.ParameterSize)
		fmt.Printf("  Quantization Level: %s\n", response.Details.QuantizationLevel)
	}

	if response.ModelInfo != nil {
		fmt.Printf("Model Info:\n")
		for key, value := range response.ModelInfo {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	if response.Parameters != "" {
		fmt.Printf("Parameters: %s\n", response.Parameters)
	}

	if response.Capabilities != nil {
		fmt.Printf("Capabilities: %v\n", response.Capabilities)
	}
}
