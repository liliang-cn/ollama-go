package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liliang-cn/ollama-go"
)

func main() {
	ctx := context.Background()

	// Pull a model with progress tracking using streaming
	fmt.Println("Pulling gemma3 model...")
	
	responseChan, errorChan := ollama.PullStream(ctx, "gemma3")

	var currentDigest string
	bars := make(map[string]*progressBar)

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// Channel closed, pull completed
				fmt.Println("\nPull completed!")
				// Clean up any remaining progress bars
				for _, bar := range bars {
					bar.Finish()
				}
				return
			}

			digest := response.Digest

			// If digest has changed and we have a bar for the previous digest, "close" it
			if digest != currentDigest && currentDigest != "" {
				if bar, exists := bars[currentDigest]; exists {
					bar.Finish()
				}
			}

			// If there's no digest, just print the status
			if digest == "" {
				fmt.Println(response.Status)
				continue
			}

			// Create a new progress bar if we don't have one for this digest and we have a total
			if _, exists := bars[digest]; !exists && response.Total > 0 {
				bars[digest] = newProgressBar(response.Total, fmt.Sprintf("pulling %s", digest[7:19]))
			}

			// Update progress bar if we have one
			if bar, exists := bars[digest]; exists && response.Completed > 0 {
				bar.Update(response.Completed)
			}

			currentDigest = digest

		case err := <-errorChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Simple progress bar implementation
type progressBar struct {
	total     int64
	current   int64
	desc      string
	finished  bool
}

func newProgressBar(total int64, desc string) *progressBar {
	return &progressBar{
		total: total,
		desc:  desc,
	}
}

func (pb *progressBar) Update(current int64) {
	if pb.finished {
		return
	}
	
	pb.current = current
	percent := float64(current) / float64(pb.total) * 100
	
	// Simple text-based progress indicator
	fmt.Printf("\r%s: %.1f%% (%s/%s)", 
		pb.desc, 
		percent, 
		formatBytes(current), 
		formatBytes(pb.total))
}

func (pb *progressBar) Finish() {
	if pb.finished {
		return
	}
	
	pb.finished = true
	fmt.Printf("\r%s: 100%% (%s/%s) ✓\n", 
		pb.desc, 
		formatBytes(pb.total), 
		formatBytes(pb.total))
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