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

	// Example digest (you would get this from CreateBlob)
	digest := "sha256:29fdb92e57cf0827ded04ae6461b5931d01fa595843f55d36f5b275a52087dd2"

	// Check if blob exists
	exists, err := client.CheckBlob(ctx, digest)
	if err != nil {
		log.Fatal("Failed to check blob:", err)
	}

	if exists {
		fmt.Printf("Blob %s exists on the server\n", digest)
	} else {
		fmt.Printf("Blob %s does not exist on the server\n", digest)
	}
}
