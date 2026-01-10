package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webhook received:", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Webhook server running on :9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return
	}
}
