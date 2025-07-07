package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func stop_auto1111() {
	fmt.Println("Stopping auto1111... (from stop_auto1111.go)")
	a1_url := "http://localhost:7860/sdapi/v1/server-stop"
	resp, err := http.Post(a1_url, "application/json", nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection was forcibly closed") {
			fmt.Println("auto1111 is shutting down (connection closed).")
			return
		}
		fmt.Println("Error stopping auto1111:", err, "(from stop_auto1111.go)")
		os.Exit(1)
	}
	defer resp.Body.Close()

	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	responseText := string(body[:n])

	if resp.StatusCode == 200 && responseText == "Stopping." {
		fmt.Println("auto1111 stopped successfully. (from stop_auto1111.go)")
	}
}
