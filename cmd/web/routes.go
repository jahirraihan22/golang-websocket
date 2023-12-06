package main

import (
	"github.com/bmizerany/pat"
	"github.com/jahirraihan22/chat/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
