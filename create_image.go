package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

func create_image(payload map[string]interface{}) bool {
    fmt.Println("Generating image...")

    // Marshal payload to JSON
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        fmt.Println("Failed to encode payload:", err)
        return false
    }

    // Send POST request
    resp, err := http.Post("http://localhost:7860/sdapi/v1/txt2img", "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        fmt.Println("Failed to call txt2img API:", err)
        return false
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        fmt.Printf("API request failed with status code: %d\n", resp.StatusCode)
        return false
    }

    // Parse response JSON
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Failed to read response:", err)
        return false
    }
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err != nil {
        fmt.Println("Failed to parse response JSON:", err)
        return false
    }

    images, ok := data["images"].([]interface{})
    if !ok || len(images) == 0 {
        fmt.Println("No image data found in response.")
        return false
    }

    // Decode base64 image
    imgBase64, ok := images[0].(string)
    if !ok {
        fmt.Println("Image data is not a string.")
        return false
    }
    imgData, err := base64.StdEncoding.DecodeString(imgBase64)
    if err != nil {
        fmt.Println("Failed to decode image:", err)
        return false
    }

    // Ask user for folder
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the output folder for the image (leave blank for current directory): ")
    userFolder, _ := reader.ReadString('\n')
    userFolder = strings.TrimSpace(userFolder)

    timestamp := time.Now().Format("150405020106") // same as %S%M%H%d%m%y
    filename := fmt.Sprintf("output_%s.png", timestamp)
    var outputPath string
    if userFolder != "" {
        outputPath = filepath.Join(userFolder, filename)
    } else {
        cwd, _ := os.Getwd()
        outputPath = filepath.Join(cwd, filename)
    }

    // Save image
    if err := os.WriteFile(outputPath, imgData, 0644); err != nil {
        fmt.Println("Failed to save image:", err)
        return false
    }
    fmt.Printf("Image saved to %s\n", outputPath)
	return true
}