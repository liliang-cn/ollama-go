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

	// Get the Ollama server version
	ctx := context.Background()
	version, err := client.Version(ctx)
	if err != nil {
		log.Fatal("Failed to get version:", err)
	}

	fmt.Printf("Ollama Server Version: %s\n", version.Version)
}