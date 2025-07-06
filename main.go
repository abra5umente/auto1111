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
	sampler := os.Getenv("SAMPLER")
	width := os.Getenv("WIDTH")
	height := os.Getenv("HEIGHT")
	steps := os.Getenv("STEPS")
	fmt.Println("Using the following settings from .env file:")
	fmt.Println("Sampler:", sampler)
	fmt.Println("Width:", width)
	fmt.Println("Height:", height)
	fmt.Println("Steps:", steps)

	// start the AUTOMATIC1111 server
	start_auto1111()
}
