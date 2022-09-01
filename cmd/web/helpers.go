package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (liberator *liberator) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	liberator.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (liberator *liberator) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (liberator *liberator) notFound(w http.ResponseWriter) {
	liberator.clientError(w, http.StatusNotFound)
}

func (liberator *liberator) logRequest(method string, code int, url string) {
	liberator.infoLog.Printf(
		"%s %v %s",
		method,
		code,
		url,
	)
}

func (liberator *liberator) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := liberator.templateCache[page]
	if !ok {
		liberator.serverError(w, fmt.Errorf("The template '%s' does not exist", page))
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		liberator.serverError(w, err)
	}
}
