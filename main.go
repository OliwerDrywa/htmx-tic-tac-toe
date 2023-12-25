package main

import (
	"fmt"
	"net/http"

	"hackathon23/handlers"
)

func main() {
	http.HandleFunc("/static/", handlers.StaticFileHandler)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/chatroom", handlers.WebSocketHandler)

	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
