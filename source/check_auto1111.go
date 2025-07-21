package main

import (
	"fmt"
	"net/http"
)

func check_auto1111() bool {
	a1_url := "http://localhost:7860/login_check/"
	resp, err := http.Get(a1_url)
	if err != nil {
		return false

	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("auto1111 is running.")
		return true
	}
	return false
}
