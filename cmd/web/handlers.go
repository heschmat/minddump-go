package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type TemplateData struct {
	CurrentYear int
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		CurrentYear: time.Now().Year(),
	}

	w.Header().Add("Server", "Go")

	// initialize a slice containing the paths to the template files
	// +++ order matters when templates share names:
	// If two files define the same template name, the last one parsed wins. +++
	tmpl_files := []string{
		"./ui/html/layout.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(tmpl_files...)
	if err != nil {
		// log.Println(err.Error())
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// if all went well, now we have:
	// Template set (ts)
	// ├── "base"
	// ├── "title"
	// └── "main"

	// err = ts.Execute(w, nil)
	err = ts.ExecuteTemplate(w, "base", data) // remember {{define "base"}} ?
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) getSnippetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "snippet %d...", id)
}
