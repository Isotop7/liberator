package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (liberator *liberator) routes() http.Handler {
	// Setup mux
	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Index handler
	mux.HandleFunc("/", liberator.dashboard)
	mux.HandleFunc("/dashboard", liberator.dashboard)

	// Books
	mux.HandleFunc("/book/create", liberator.bookCreate)
	mux.HandleFunc("/book/view", liberator.bookView)

	// Default middleware chain
	defaultChain := alice.New(liberator.recoverPanic, liberator.logRequest, secureHeaders)

	return defaultChain.Then(mux)
}
