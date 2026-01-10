package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
	body, _ := io.ReadAll(r.Body)
	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Обработчик всех путей
	http.HandleFunc("/", handler)

	fmt.Println("Webhook server running on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
