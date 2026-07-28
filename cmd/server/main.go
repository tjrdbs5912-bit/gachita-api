package main

import (
	"fmt"
	"net/http"

	httpadapter "gachita-api/internal/adapter/http"
)

func main() {
	handler := httpadapter.NewRouter()

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
