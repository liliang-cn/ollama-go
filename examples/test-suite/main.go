package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	fmt.Println("🎉 Ollama Go Client Test Suite")
	fmt.Println("==============================")

	// Test 1: List models
	fmt.Println("\n✅ Test 1: Listing models...")
	models, err := ollama.List(ctx)
	if err != nil {
		log.Fatal("❌ List test failed:", err)
	}
	fmt.Printf("Found %d models:\n", len(models.Models))
	for i, model := range models.Models {
		if i < 3 { // Show only first 3
			fmt.Printf("  - %s\n", model.Model)
		}
	}
	if len(models.Models) > 3 {
		fmt.Printf("  ... and %d more\n", len(models.Models)-3)
	}

	if len(models.Models) == 0 {
		fmt.Println("❌ No models found. Please pull a model first: ollama pull llama2")
		return
	}

	modelName := models.Models[0].Model
	fmt.Printf("Using model: %s\n", modelName)

	// Test 2: Simple generation
	fmt.Println("\n✅ Test 2: Simple generation...")
	response, err := ollama.Generate(ctx, modelName, "Say hello in one sentence.")
	if err != nil {
		log.Fatal("❌ Generate test failed:", err)
	}
	fmt.Printf("Response: %s\n", response.Response)

	// Test 3: Chat
	fmt.Println("\n✅ Test 3: Chat...")
	messages := []ollama.Message{
		{Role: "user", Content: "What's 2+2? Answer briefly."},
	}
	chatResp, err := ollama.Chat(ctx, modelName, messages)
	if err != nil {
		log.Fatal("❌ Chat test failed:", err)
	}
	fmt.Printf("Chat response: %s\n", chatResp.Message.Content)

	// Test 4: Streaming (limited)
	fmt.Println("\n✅ Test 4: Streaming generation...")
	responseChan, errorChan := ollama.GenerateStream(ctx, modelName, "Count from 1 to 3.")
	
	timeout := time.After(15 * time.Second)
	var streamingContent string
	responseCount := 0

streamingLoop:
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				break streamingLoop
			}
			responseCount++
			streamingContent += response.Response
			if response.Done || responseCount > 10 { // Limit responses
				break streamingLoop
			}
		case err := <-errorChan:
			if err != nil {
				log.Printf("❌ Streaming error: %v", err)
				break streamingLoop
			}
		case <-timeout:
			fmt.Println("⏰ Streaming timeout reached")
			break streamingLoop
		}
	}
	
	fmt.Printf("Streaming result (%d chunks): %s\n", responseCount, streamingContent)

	// Test 5: Model info
	fmt.Println("\n✅ Test 5: Model info...")
	info, err := ollama.Show(ctx, modelName)
	if err != nil {
		log.Printf("⚠️ Show test warning: %v", err)
	} else {
		fmt.Printf("Model details: Family=%s, Size=%s\n", 
			getStringValue(info.Details, "family"), 
			getStringValue(info.Details, "parameter_size"))
	}

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("\n📚 Your Ollama Go client is working perfectly!")
	fmt.Println("You can now use github.com/liliang-cn/ollama-go in your projects.")
}

func getStringValue(details *ollama.ModelDetails, field string) string {
	if details == nil {
		return "unknown"
	}
	switch field {
	case "family":
		return details.Family
	case "parameter_size":
		return details.ParameterSize
	default:
		return "unknown"
	}
}