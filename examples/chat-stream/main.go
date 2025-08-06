package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	messages := []ollama.Message{
		{
			Role:    "user", 
			Content: "Tell me about Go programming language",
		},
	}

	// Stream chat
	responseChan, errorChan := ollama.ChatStream(ctx, "gemma3", messages)

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				fmt.Println("\nDone!")
				return
			}
			fmt.Print(response.Message.Content)
		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}