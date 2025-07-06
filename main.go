package main

import (
	// "bufio"          // for reading user input
	// "bytes"          // for constructing POST payloads
	// "encoding/base64" // decoding base64 image data
	// "encoding/json"  // marshaling/unmarshaling JSON
	"fmt" // printing to console
	//  "io"             // for reading HTTP responses
	//  "net/http"       // making HTTP requests
	"os" // environment variables, file operations
	//  "path/filepath"  // for building save paths
	// "strconv"        // string to int/float conversions
	// "strings"        // string manipulation
	//  "time"           // sleeping, timestamps
	"github.com/joho/godotenv" // loading environment variables from .env file
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using system environment variables.")
	}
	sampler := os.Getenv("SAMPLER_NAME")
	scheduler := os.Getenv("SCHEDULER_NAME")
	width := os.Getenv("IMAGE_WIDTH")
	height := os.Getenv("IMAGE_HEIGHT")
	steps := os.Getenv("STEPS")
	fmt.Println("Using the following settings from .env file:")
	fmt.Println("Sampler:", sampler)
	fmt.Println("Scheduler:", scheduler)
	fmt.Println("Steps:", steps)
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)

	// start the AUTOMATIC1111 server
	start_auto1111()
}
