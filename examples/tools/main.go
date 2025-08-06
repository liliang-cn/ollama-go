package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

// Tool functions
func addTwoNumbers(a, b int) int {
	return a + b
}

func subtractTwoNumbers(a, b int) int {
	return a - b
}

func main() {
	ctx := context.Background()

	// Define tools - one using function definition, one manual
	addTool := &ollama.Tool{
		Type: "function",
		Function: &ollama.ToolFunction{
			Name:        "add_two_numbers",
			Description: "Add two numbers",
			Parameters: map[string]interface{}{
				"type":     "object",
				"required": []string{"a", "b"},
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "integer",
						"description": "The first number",
					},
					"b": map[string]interface{}{
						"type":        "integer",
						"description": "The second number",
					},
				},
			},
		},
	}

	// Manual tool definition (like Python example)
	subtractTool := &ollama.Tool{
		Type: "function",
		Function: &ollama.ToolFunction{
			Name:        "subtract_two_numbers",
			Description: "Subtract two numbers",
			Parameters: map[string]interface{}{
				"type":     "object",
				"required": []string{"a", "b"},
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type":        "integer",
						"description": "The first number",
					},
					"b": map[string]interface{}{
						"type":        "integer",
						"description": "The second number",
					},
				},
			},
		},
	}

	messages := []ollama.Message{
		{Role: "user", Content: "What is three plus one?"},
	}
	fmt.Printf("Prompt: %s\n", messages[0].Content)

	// Available functions map (like Python)
	availableFunctions := map[string]func(int, int) int{
		"add_two_numbers":      addTwoNumbers,
		"subtract_two_numbers": subtractTwoNumbers,
	}

	// Make initial chat request with tools
	response, err := ollama.Chat(ctx, "llama3.1", messages, func(req *ollama.ChatRequest) {
		req.Tools = []ollama.Tool{*addTool, *subtractTool}
	})
	if err != nil {
		log.Fatal(err)
	}

	var output int
	
	if len(response.Message.ToolCalls) > 0 {
		// Process tool calls (may be multiple)
		for _, toolCall := range response.Message.ToolCalls {
			// Check if function is available
			if function, exists := availableFunctions[toolCall.Function.Name]; exists {
				fmt.Printf("Calling function: %s\n", toolCall.Function.Name)
				fmt.Printf("Arguments: %v\n", toolCall.Function.Arguments)
				
				// Extract arguments and call function
				a := int(toolCall.Function.Arguments["a"].(float64))
				b := int(toolCall.Function.Arguments["b"].(float64))
				output = function(a, b)
				
				fmt.Printf("Function output: %d\n", output)
			} else {
				fmt.Printf("Function %s not found\n", toolCall.Function.Name)
			}
		}

		// Add the tool call result to messages
		messages = append(messages, response.Message)
		messages = append(messages, ollama.Message{
			Role:     "tool",
			Content:  fmt.Sprintf("%d", output),
			ToolName: response.Message.ToolCalls[0].Function.Name,
		})

		// Get final response with tool results
		finalResponse, err := ollama.Chat(ctx, "llama3.1", messages)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Final response: %s\n", finalResponse.Message.Content)
	} else {
		fmt.Println("No tool calls returned from model")
	}
}