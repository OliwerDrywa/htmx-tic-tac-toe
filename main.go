package main

import (
	"fmt"
	"net/http"

	"hackathon23/handlers"
)

func main() {
	http.HandleFunc("/static/", handlers.StaticFileHandler)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/tic-tac-toe", handlers.TicTacToeHandler)
	http.HandleFunc("/socket", handlers.WebSocketHandler)

	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
