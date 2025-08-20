package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "翻译为英语：我们后面是不是也需要为GUI的VSAN mode做设计？页面的差别挺大的，HCI的只有一个页面和GUI不一样，应该比较简单。",
		},
	}

	// Stream chat
	responseChan, errorChan := ollama.ChatStream(ctx, "qwen3", messages)

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				fmt.Println("\nDone!")
				return
			}
			fmt.Print(response.Message.Content)
		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
