package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

// Define the schema for the response using Go structs
type FriendInfo struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	IsAvailable bool   `json:"is_available"`
}

type FriendList struct {
	Friends []FriendInfo `json:"friends"`
}

func main() {
	ctx := context.Background()

	// Create the JSON schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"friends": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"age": map[string]interface{}{
							"type": "integer",
						},
						"is_available": map[string]interface{}{
							"type": "boolean",
						},
					},
					"required": []string{"name", "age", "is_available"},
				},
			},
		},
		"required": []string{"friends"},
	}

	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "I have two friends. The first is Ollama 22 years old busy saving the world, and the second is Alonso 23 years old and wants to hang out. Return a list of friends in JSON format",
		},
	}

	// Make request with structured output
	response, err := ollama.Chat(ctx, "llama3.1:8b", messages, func(req *ollama.ChatRequest) {
		req.Format = schema
		req.Options = &ollama.Options{
			Temperature: ollama.Float64Ptr(0), // Make responses more deterministic
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Parse and validate the response
	var friendsList FriendList
	err = json.Unmarshal([]byte(response.Message.Content), &friendsList)
	if err != nil {
		log.Fatal("Failed to parse response:", err)
	}

	fmt.Printf("%+v\n", friendsList)
}