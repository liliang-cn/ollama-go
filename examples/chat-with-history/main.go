package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	// Initial conversation history
	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "Why is the sky blue?",
		},
		{
			Role:    "assistant",
			Content: "The sky is blue because of the way the Earth's atmosphere scatters sunlight.",
		},
		{
			Role:    "user",
			Content: "What is the weather in Tokyo?",
		},
		{
			Role:    "assistant",
			Content: "The weather in Tokyo is typically warm and humid during the summer months, with temperatures often exceeding 30°C (86°F). The city experiences a rainy season from June to September, with heavy rainfall and occasional typhoons. Winter is mild, with temperatures rarely dropping below freezing. The city is known for its high-tech and vibrant culture, with many popular tourist attractions such as the Tokyo Tower, Senso-ji Temple, and the bustling Shibuya district.",
		},
	}

	for {
		fmt.Print("Chat with history: ")
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()
		if userInput == "" {
			continue
		}

		// Add user input to messages
		currentMessages := append(messages, ollama.Message{
			Role:    "user",
			Content: userInput,
		})

		response, err := ollama.Chat(ctx, "gemma3", currentMessages)
		if err != nil {
			log.Fatal(err)
		}

		// Add the response to the messages to maintain the history
		messages = append(messages, 
			ollama.Message{Role: "user", Content: userInput},
			ollama.Message{Role: "assistant", Content: response.Message.Content},
		)

		fmt.Printf("%s\n\n", response.Message.Content)
	}
}