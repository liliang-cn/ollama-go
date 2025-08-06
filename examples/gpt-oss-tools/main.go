package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/liliang-cn/ollama-go"
)

func getWeather(city string) string {
	temperatures := []int{}
	for i := -10; i <= 35; i++ {
		temperatures = append(temperatures, i)
	}
	temp := temperatures[rand.Intn(len(temperatures))]
	return fmt.Sprintf("The temperature in %s is %d°C", city, temp)
}

func getWeatherConditions(city string) string {
	conditions := []string{"sunny", "cloudy", "rainy", "snowy", "foggy"}
	return conditions[rand.Intn(len(conditions))]
}

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	// Define available tools
	weatherTool := ollama.Tool{
		Type: "function",
		Function: &ollama.ToolFunction{
			Name:        "get_weather",
			Description: "Get the current temperature for a city",
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
			Name:        "get_weather_conditions",
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

	messages := []ollama.Message{
		{Role: "user", Content: "What is the weather like in London? What are the conditions in Toronto?"},
	}

	model := "gpt-oss:20b"

	// gpt-oss can call tools while "thinking" 
	// a loop is needed to call the tools and get the results
	for {
		response, err := ollama.Chat(ctx, model, messages, ollama.WithTools([]ollama.Tool{weatherTool, conditionsTool}))
		if err != nil {
			log.Fatal(err)
		}

		if response.Message.Content != "" {
			fmt.Println("Content:")
			fmt.Printf("%s\n\n", response.Message.Content)
		}
		if response.Message.Thinking != "" {
			fmt.Println("Thinking:")
			fmt.Printf("%s\n\n", response.Message.Thinking)
		}

		if len(response.Message.ToolCalls) > 0 {
			for _, toolCall := range response.Message.ToolCalls {
				var result string
				switch toolCall.Function.Name {
				case "get_weather":
					city := toolCall.Function.Arguments["city"].(string)
					result = getWeather(city)
				case "get_weather_conditions":
					city := toolCall.Function.Arguments["city"].(string)
					result = getWeatherConditions(city)
				default:
					fmt.Printf("Tool %s not found\n", toolCall.Function.Name)
					continue
				}

				fmt.Printf("Result from tool call name: %s with arguments: %v result: %s\n\n", 
					toolCall.Function.Name, toolCall.Function.Arguments, result)

				// Add the assistant message and tool result to the messages
				messages = append(messages, response.Message)
				messages = append(messages, ollama.Message{
					Role:     "tool",
					Content:  result,
					ToolName: toolCall.Function.Name,
				})
			}
		} else {
			// no more tool calls, we can stop the loop
			break
		}
	}
}