package main

import (
	"github.com/jahirraihan22/chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("Starting Channel listener .. ")
	go handlers.ListenToWsChannel()

	log.Println("Server is starting...")
	_ = http.ListenAndServe(":8080", mux)
}
