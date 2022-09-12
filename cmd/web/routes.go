package main

import (
	"net/http"

	"github.com/Isotop7/liberator/ui"
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
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Session agnostic handler
	publicCalls := alice.New(liberator.sessionManager.LoadAndSave, secureCSRF, liberator.authenticate)

	// Serve handlers
	router.Handler(http.MethodGet, "/", publicCalls.ThenFunc(liberator.dashboard))
	router.Handler(http.MethodGet, "/dashboard", publicCalls.ThenFunc(liberator.dashboard))
	router.Handler(http.MethodGet, "/user/signup", publicCalls.ThenFunc(liberator.userSignup))
	router.Handler(http.MethodPost, "/user/signup", publicCalls.ThenFunc(liberator.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", publicCalls.ThenFunc(liberator.userLogin))
	router.Handler(http.MethodPost, "/user/login", publicCalls.ThenFunc(liberator.userLoginPost))

	protectedCalls := publicCalls.Append(liberator.requireAuthentication)
	router.Handler(http.MethodGet, "/book/view/:id", protectedCalls.ThenFunc(liberator.bookView))
	router.Handler(http.MethodPost, "/search", protectedCalls.ThenFunc(liberator.searchView))
	router.Handler(http.MethodGet, "/book/create", protectedCalls.ThenFunc(liberator.bookCreate))
	router.Handler(http.MethodPost, "/book/create", protectedCalls.ThenFunc(liberator.bookCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protectedCalls.ThenFunc(liberator.userLogoutPost))

	// Default middleware chain
	standard := alice.New(liberator.recoverPanic, liberator.logRequest, secureHeaders)

	return standard.Then(router)
}
