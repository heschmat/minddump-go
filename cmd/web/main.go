package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/")) // &{./ui/static/}
	// N.B. actual file access happens per request.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", getSnippetByID)

	log.Println("starting the server on :4000")

	// params: the TCP network address, the servemux
	err := http.ListenAndServe(":4000", mux)
	// N.B. Any error returned by http.ListenAndServe() is always non-nil.
	log.Fatal(err)
}
