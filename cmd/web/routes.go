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

	// Session agnostic handler
	dynamic := alice.New(liberator.sessionManager.LoadAndSave)

	// Serve handlers
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(liberator.dashboard))
	router.Handler(http.MethodGet, "/dashboard", dynamic.ThenFunc(liberator.dashboard))
	router.Handler(http.MethodGet, "/book/view/:id", dynamic.ThenFunc(liberator.bookView))
	router.Handler(http.MethodGet, "/book/create", dynamic.ThenFunc(liberator.bookCreate))
	router.Handler(http.MethodPost, "/book/create", dynamic.ThenFunc(liberator.bookCreatePost))
	router.Handler(http.MethodPost, "/search", dynamic.ThenFunc(liberator.searchView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(liberator.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(liberator.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(liberator.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(liberator.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(liberator.userLogoutPost))

	// Default middleware chain
	standard := alice.New(liberator.recoverPanic, liberator.logRequest, secureHeaders)

	return standard.Then(router)
}
