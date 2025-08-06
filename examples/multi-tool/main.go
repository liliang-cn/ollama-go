package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/liliang-cn/ollama-go"
)

// Mock weather functions
func getTemperature(city string) string {
	validCities := []string{"London", "Paris", "New York", "Tokyo", "Sydney"}
	for _, validCity := range validCities {
		if city == validCity {
			temp := rand.Intn(35) // 0-35 degrees
			return fmt.Sprintf("%d degrees Celsius", temp)
		}
	}
	return "Unknown city"
}

func getConditions(city string) string {
	validCities := []string{"London", "Paris", "New York", "Tokyo", "Sydney"}
	for _, validCity := range validCities {
		if city == validCity {
			conditions := []string{"sunny", "cloudy", "rainy", "snowy"}
			return conditions[rand.Intn(len(conditions))]
		}
	}
	return "Unknown city"
}

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// Define tools
	tempTool := ollama.Tool{
		Type: "function",
		Function: &ollama.ToolFunction{
			Name:        "get_temperature",
			Description: "Get the temperature for a city in Celsius",
			Parameters: map[string]interface{}{
				"type": "object",
				"required": []string{"city"},
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "The name of the city",
					},
				},
			},
		},
	}

	conditionsTool := ollama.Tool{
		Type: "function",
		Function: &ollama.ToolFunction{
			Name:        "get_conditions",
			Description: "Get the weather conditions for a city",
			Parameters: map[string]interface{}{
				"type": "object",
				"required": []string{"city"},
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "The name of the city",
					},
				},
			},
		},
	}

	cities := []string{"London", "Paris", "New York", "Tokyo", "Sydney"}
	city1 := cities[rand.Intn(len(cities))]
	city2 := cities[rand.Intn(len(cities))]

	messages := []ollama.Message{
		{Role: "user", Content: fmt.Sprintf("What is the temperature in %s? and what are the weather conditions in %s?", city1, city2)},
	}
	
	fmt.Printf("----- Prompt: %s\n\n", messages[0].Content)

	// Make initial chat request with streaming and thinking
	options := func(req *ollama.ChatRequest) {
		req.Stream = ollama.BoolPtr(true)
		req.Tools = []ollama.Tool{tempTool, conditionsTool}
		req.Think = ollama.BoolPtr(true)
	}

	responsesCh, errorsCh := ollama.ChatStream(ctx, "qwen3", messages, options)

	for {
		select {
		case response, ok := <-responsesCh:
			if !ok {
				goto processToolCalls
			}
			// Handle thinking output
			if response.Message.Thinking != "" {
				fmt.Print(response.Message.Thinking)
			}
			// Handle regular content
			if response.Message.Content != "" {
				fmt.Print(response.Message.Content)
			}
			// Handle tool calls
			if len(response.Message.ToolCalls) > 0 {
				for _, toolCall := range response.Message.ToolCalls {
					var output string
					switch toolCall.Function.Name {
					case "get_temperature":
						city := toolCall.Function.Arguments["city"].(string)
						fmt.Printf("\nCalling function: %s with arguments: %v\n", toolCall.Function.Name, toolCall.Function.Arguments)
						output = getTemperature(city)
						fmt.Printf("> Function output: %s\n\n", output)
					case "get_conditions":
						city := toolCall.Function.Arguments["city"].(string)
						fmt.Printf("\nCalling function: %s with arguments: %v\n", toolCall.Function.Name, toolCall.Function.Arguments)
						output = getConditions(city)
						fmt.Printf("> Function output: %s\n\n", output)
					}

					// Add tool call and result to messages
					messages = append(messages, response.Message)
					messages = append(messages, ollama.Message{
						Role:    "tool",
						Content: output,
					})
				}
			}
		case err := <-errorsCh:
			if err != nil {
				log.Fatal(err)
			}
		}
	}

processToolCalls:
	fmt.Println("----- Sending result back to model\n")

	// Check if we have tool results to send back
	hasToolResults := false
	for _, msg := range messages {
		if msg.Role == "tool" {
			hasToolResults = true
			break
		}
	}

	if hasToolResults {
		// Send final request with tool results
		finalOptions := func(req *ollama.ChatRequest) {
			req.Stream = ollama.BoolPtr(true)
			req.Tools = []ollama.Tool{tempTool, conditionsTool}
			req.Think = ollama.BoolPtr(true)
		}

		finalResponsesCh, finalErrorsCh := ollama.ChatStream(ctx, "qwen3", messages, finalOptions)
		doneThinking := false

		for {
			select {
			case response, ok := <-finalResponsesCh:
				if !ok {
					return
				}
				if response.Message.Thinking != "" {
					fmt.Print(response.Message.Thinking)
				}
				if response.Message.Content != "" {
					if !doneThinking {
						fmt.Println("\n----- Final result:")
						doneThinking = true
					}
					fmt.Print(response.Message.Content)
				}
				if len(response.Message.ToolCalls) > 0 {
					// Model should be explaining the tool calls and the results in this output
					fmt.Println("Model returned tool calls:")
					fmt.Printf("%+v\n", response.Message.ToolCalls)
				}
			case err := <-finalErrorsCh:
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		fmt.Println("No tool calls returned")
	}
}