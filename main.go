package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/images/:id", handleDicomImage)

	err := http.ListenAndServe(":3333", server)
	if err != nil {
		fmt.Printf("Error while opening the server: %v", err)
	}
}

func handleDicomImage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Get an image
	case http.MethodPost:
		// TODO: Store a new image
	default:
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}
