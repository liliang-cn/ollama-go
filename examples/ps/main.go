package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// List running processes
	processes, err := ollama.Ps(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running Models:")
	fmt.Printf("%-20s %-10s %-10s %-20s\n", "Model", "Size", "VRAM", "Expires At")
	fmt.Println(strings.Repeat("-", 70))

	for _, model := range processes.Models {
		sizeStr := "N/A"
		if model.Size > 0 {
			sizeStr = formatBytes(model.Size)
		}

		vramStr := "N/A"
		if model.SizeVRAM > 0 {
			vramStr = formatBytes(model.SizeVRAM)
		}

		expiresStr := "N/A"
		if model.ExpiresAt != nil {
			expiresStr = model.ExpiresAt.Format("15:04:05")
		}

		fmt.Printf("%-20s %-10s %-10s %-20s\n",
			model.Name,
			sizeStr,
			vramStr,
			expiresStr)
	}

	if len(processes.Models) == 0 {
		fmt.Println("No models are currently running")
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
