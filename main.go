package main

import (
	"net/http"

	"hackathon23/handlers"
)

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.ListenAndServe(":8080", nil)
}
