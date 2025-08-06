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
			Content: "Why is the sky blue?",
		},
	}

	response, err := ollama.Chat(ctx, "gemma3", messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Message.Content)
}
