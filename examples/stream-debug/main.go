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

	// Test with streaming to see if we can get more detailed info
	fmt.Println("Testing streaming generation...")
	
	responseChan, errorChan := ollama.GenerateStream(ctx, "gpt-oss:20b", "Hello! Please say hi back.")

	timeout := time.After(10 * time.Second)
	var fullResponse string
	responseCount := 0

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				fmt.Printf("\nResponse channel closed. Total responses: %d\n", responseCount)
				fmt.Printf("Full response: '%s'\n", fullResponse)
				return
			}
			responseCount++
			fmt.Printf("Response %d: Done=%t, Content='%s'\n", responseCount, response.Done, response.Response)
			fullResponse += response.Response
			
			if response.Done {
				fmt.Printf("\nGeneration complete! Full response: '%s'\n", fullResponse)
				return
			}

		case err := <-errorChan:
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}

		case <-timeout:
			fmt.Println("\nTimeout reached, no response within 10 seconds")
			return
		}
	}
}