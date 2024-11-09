package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, custom string, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("custom", custom), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	http.Error(w, fmt.Sprintf("Issue with: %s with error: %s", custom, err.Error()), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int, custom string) {
	http.Error(w, fmt.Sprintf("with error: %s", custom), status)
}
