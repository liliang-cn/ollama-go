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

	// Get image path from user
	fmt.Print("Please enter the path to the image: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	// Create message with image
	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "What is in this image? Be concise.",
			Images:  []ollama.Image{{Data: imagePath}},
		},
	}

	// Send chat request
	response, err := ollama.Chat(ctx, "gemma3", messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Message.Content)
}
