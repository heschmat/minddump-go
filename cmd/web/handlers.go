package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/heschmat/minddump-go/internal/models"
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
		"./ui/html/partials/nav.tmpl.htmlx",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(tmpl_files...)
	if err != nil {
		// log.Println(err.Error())
		app.error.serverError(w, r, err)
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
		// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.error.serverError(w, r, err)
		return
	}
}

func (app *application) getSnippetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.error.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	// write the snippet data as a plain-text HTTP response body for now, to test the Get() method of SnippetModel
	fmt.Fprintf(w, "%+v", snippet)
	// fmt.Fprintf(w, "snippet %d...", id)
}

func (app *application) postSnippet(w http.ResponseWriter, r *http.Request) {
	// dummy data to test the Insert() method of SnippetModel
	title := "Security Practice"
	content := "Rotate secrets regularly and never commit .env files to the repository."
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.error.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
