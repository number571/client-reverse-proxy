package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	for {
		time.Sleep(time.Second)
		req, err := http.NewRequest(http.MethodGet, "http://gateway:9000/ping", nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		httpClient := &http.Client{Timeout: 5 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			fmt.Println(err)
			continue
		}
		resp.Body.Close()

		fmt.Println(string(body))
	}
}
