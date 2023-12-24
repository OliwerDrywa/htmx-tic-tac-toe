package main

import (
	"fmt"
	"net/http"

	"hackathon23/handlers"
)

func main() {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/chatroom", handlers.WebSocketHandler)

	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
