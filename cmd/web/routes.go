package main

import "net/http"

func (liberator *liberator) routes() *http.ServeMux {
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

	return mux
}
