package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func stop_auto1111(cmd *exec.Cmd) {
	fmt.Println("Stopping auto1111...")
	a1_url := "http://localhost:7860/sdapi/v1/server-stop"
	_, err := http.Post(a1_url, "application/json", nil)
	if err != nil {
		if strings.Contains(err.Error(), "connection was forcibly closed") {
			fmt.Println("auto1111 is shutting down (connection closed).")
			return
		}
		fmt.Println("Error stopping auto1111 via API:", err)
		fmt.Println("Attempting to terminate the process directly...")
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Failed to kill process:", err)
		} else {
			fmt.Println("Process terminated.")
		}
		return
	}

	fmt.Println("Waiting for auto1111 to stop...")
	for i := 0; i < 10; i++ {
		if !check_auto1111() {
			fmt.Println("auto1111 stopped successfully.")
			return
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("auto1111 did not stop gracefully. Terminating process...")
	if err := cmd.Process.Kill(); err != nil {
		fmt.Println("Failed to kill process:", err)
	} else {
		fmt.Println("Process terminated.")
	}
}
