package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	response, err := ollama.Generate(ctx, "deepseek-r1", "why is the sky blue", func(req *ollama.GenerateRequest) {
		req.Think = ollama.BoolPtr(true)
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Thinking:\n========\n")
	fmt.Println(response.Thinking)
	fmt.Println("\nResponse:\n========\n")
	fmt.Println(response.Response)
}