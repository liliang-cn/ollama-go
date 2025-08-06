package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Create a custom model with streaming progress
	modelfile := `FROM gemma3
SYSTEM "You are mario from Super Mario Bros."`

	progressChan, errorChan := ollama.CreateStream(ctx, "my-assistant", modelfile)

	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				fmt.Println("Model created successfully!")
				return
			}
			fmt.Printf("Status: %s\n", progress.Status)
		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}