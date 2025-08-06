package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/liliang-cn/ollama-go"
)

// Define the schema for image objects
type Object struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	Attributes string  `json:"attributes"`
}

type ImageDescription struct {
	Summary     string   `json:"summary"`
	Objects     []Object `json:"objects"`
	Scene       string   `json:"scene"`
	Colors      []string `json:"colors"`
	TimeOfDay   string   `json:"time_of_day"`   // Morning, Afternoon, Evening, Night
	Setting     string   `json:"setting"`       // Indoor, Outdoor, Unknown
	TextContent *string  `json:"text_content"`
}

func main() {
	ctx := context.Background()

	// Get path from user input
	fmt.Print("Enter the path to your image: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	// Verify the file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image not found at: %s", imagePath)
	}

	// Create the JSON schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"summary": map[string]interface{}{
				"type": "string",
			},
			"objects": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"confidence": map[string]interface{}{
							"type": "number",
						},
						"attributes": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"name", "confidence", "attributes"},
				},
			},
			"scene": map[string]interface{}{
				"type": "string",
			},
			"colors": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"time_of_day": map[string]interface{}{
				"type": "string",
				"enum": []string{"Morning", "Afternoon", "Evening", "Night"},
			},
			"setting": map[string]interface{}{
				"type": "string",
				"enum": []string{"Indoor", "Outdoor", "Unknown"},
			},
			"text_content": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"summary", "objects", "scene", "colors", "time_of_day", "setting"},
	}

	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "Analyze this image and return a detailed JSON description including objects, scene, colors and any text detected. If you cannot determine certain details, leave those fields empty.",
			Images:  []ollama.Image{{Data: imagePath}},
		},
	}

	// Set up chat with structured output
	response, err := ollama.Chat(ctx, "gemma3", messages, func(req *ollama.ChatRequest) {
		req.Format = schema
		req.Options = &ollama.Options{
			Temperature: ollama.Float64Ptr(0), // Set temperature to 0 for more deterministic output
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	// Convert received content to the schema
	var imageAnalysis ImageDescription
	err = json.Unmarshal([]byte(response.Message.Content), &imageAnalysis)
	if err != nil {
		log.Fatal("Failed to parse response:", err)
	}

	fmt.Printf("%+v\n", imageAnalysis)
}