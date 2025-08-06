package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Stream generation
	responseChan, errorChan := ollama.GenerateStream(ctx, "gemma3", "Tell me a story about a brave little mouse.")

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				fmt.Println("\nDone!")
				return
			}
			fmt.Print(response.Response)
		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}