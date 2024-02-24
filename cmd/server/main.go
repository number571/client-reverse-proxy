package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s]: ping\n", r.RemoteAddr)
		fmt.Fprint(w, "pong")
	})
	fmt.Println("server is listening...")
	http.ListenAndServe(":8080", nil)
}
