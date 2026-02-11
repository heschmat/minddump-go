package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Helloooooo"))
}

func getSnippetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "snippet %d...", id)
}
