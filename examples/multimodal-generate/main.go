package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Get comic number from command line arguments, or use random
	var num int
	if len(os.Args) > 1 {
		var err error
		num, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("Invalid comic number:", err)
		}
	} else {
		// Get latest comic to find max number
		resp, err := http.Get("https://xkcd.com/info.0.json")
		if err != nil {
			log.Fatal("Failed to get latest comic info:", err)
		}
		defer resp.Body.Close()

		// For simplicity, use a random number between 1-2000
		num = 1 + (int(resp.Header.Get("Content-Length")[0]) % 2000)
	}

	// Get comic info
	comicResp, err := http.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", num))
	if err != nil {
		log.Fatal("Failed to get comic info:", err)
	}
	defer comicResp.Body.Close()

	// Download comic image
	imageResp, err := http.Get(fmt.Sprintf("https://imgs.xkcd.com/comics/%d.png", num))
	if err != nil {
		log.Fatal("Failed to download comic image:", err)
	}
	defer imageResp.Body.Close()

	// Read image data
	imageData := make([]byte, imageResp.ContentLength)
	_, err = imageResp.Body.Read(imageData)
	if err != nil {
		log.Fatal("Failed to read image data:", err)
	}

	fmt.Printf("xkcd #%d\n", num)
	fmt.Printf("link: https://xkcd.com/%d\n", num)
	fmt.Println("---")

	// Use streaming generation to explain the comic
	responsesCh, errorsCh := ollama.GenerateStream(ctx, "llava", "explain this comic:",
		func(req *ollama.GenerateRequest) {
			req.Images = []ollama.Image{{Data: string(imageData)}}
		})

	for {
		select {
		case response, ok := <-responsesCh:
			if !ok {
				fmt.Println()
				return
			}
			fmt.Print(response.Response)
		case err := <-errorsCh:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
