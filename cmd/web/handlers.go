package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), slog.String("msg", "parsing the Templates"), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		http.Error(w, fmt.Sprintf("Internal Server Error\nissue parsing the Templates with error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), slog.String("msg", "executing the Templates"), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		http.Error(w, fmt.Sprintf("Internal Server Error\n%sIssue executing the templates with error: ", err.Error()), http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
