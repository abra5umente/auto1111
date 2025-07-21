package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func create_image(payload map[string]interface{}) (string, bool) {
	var outputPath string
	if out, ok := payload["output"].(string); ok && out != "" {
		outputPath = out
	} else {
		ex, _ := os.Executable()
		projectRoot := filepath.Clean(filepath.Join(filepath.Dir(ex), ".."))
		outputDir := filepath.Join(projectRoot, "output")
		_ = os.MkdirAll(outputDir, os.ModePerm)

		timestamp := time.Now().Format("150405020106")
		outputPath = filepath.Join(outputDir, fmt.Sprintf("output-%s.png", timestamp))
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return "", false
	}

	// Send POST request to auto1111 API
	a1_url := "http://localhost:7860/sdapi/v1/txt2img"
	resp, err := http.Post(a1_url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending request to auto1111:", err)
		return "", false
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", false
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response JSON:", err)
		return "", false
	}

	// Extract and decode image
	if images, ok := result["images"].([]interface{}); ok && len(images) > 0 {
		if imgStr, ok := images[0].(string); ok {
			decodedImage, err := base64.StdEncoding.DecodeString(imgStr)
			if err != nil {
				fmt.Println("Error decoding base64 image:", err)
				return "", false
			}

			// Save image
			err = os.WriteFile(outputPath, decodedImage, 0644)
			if err != nil {
				fmt.Println("Error saving image:", err)
				return "", false
			}
			return outputPath, true
		}
	}

	fmt.Println("No image found in the response.")
	return "", false
}
