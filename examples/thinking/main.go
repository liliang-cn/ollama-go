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
			Content: "What is 10 + 23?",
		},
	}

	response, err := ollama.Chat(ctx, "deepseek-r1", messages, func(req *ollama.ChatRequest) {
		req.Think = ollama.BoolPtr(true)
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Thinking:")
	fmt.Println("========")
	fmt.Println()
	fmt.Println(response.Message.Thinking)
	fmt.Println()
	fmt.Println("Response:")
	fmt.Println("========")
	fmt.Println()
	fmt.Println(response.Message.Content)
}
