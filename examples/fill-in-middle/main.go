package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	prompt := `def remove_non_ascii(s: str) -> str:
    """ `

	suffix := `
    return result
`

	response, err := ollama.Generate(ctx, "codellama:7b-code", prompt, func(req *ollama.GenerateRequest) {
		req.Suffix = suffix
		req.Options = &ollama.Options{
			NumPredict:  ollama.IntPtr(128),
			Temperature: ollama.Float64Ptr(0),
			TopP:        ollama.Float64Ptr(0.9),
			Stop:        []string{"<EOT>"},
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Response)
}