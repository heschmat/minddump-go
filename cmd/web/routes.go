package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/")) // &{./ui/static/}
	// N.B. actual file access happens per request.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.getSnippetByID)
	mux.HandleFunc("POST /snippet/create", app.postSnippet)

	return mux
}
