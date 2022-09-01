package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (liberator *liberator) routes() http.Handler {
	// Setup router
	router := httprouter.New()

	// Custom error handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		liberator.notFound(w)
	})

	// Serve static files
	fileServer := http.FileServer(http.Dir("./assets/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Serve handlers
	router.HandlerFunc(http.MethodGet, "/", liberator.dashboard)
	router.HandlerFunc(http.MethodGet, "/dashboard", liberator.dashboard)
	router.HandlerFunc(http.MethodGet, "/book/view/:id", liberator.bookView)
	router.HandlerFunc(http.MethodGet, "/book/create", liberator.bookCreate)
	router.HandlerFunc(http.MethodPost, "/book/create", liberator.bookCreatePost)

	// Default middleware chain
	defaultChain := alice.New(liberator.recoverPanic, liberator.logRequest, secureHeaders)

	return defaultChain.Then(router)
}
