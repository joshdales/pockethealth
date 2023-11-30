package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := http.NewServeMux()

	err := http.ListenAndServe(":3333", server)
	if err != nil {
		fmt.Printf("Error while opening the server: %v", err)
	}
}
