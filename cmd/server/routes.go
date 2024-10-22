package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	// Equivalent to: mux.Handle("/", http.HandlerFunc(home))
	mux.HandleFunc("/blog/view", app.blogView)
	mux.HandleFunc("/blog/create", app.blogCreate)

	return secureHeaders(mux)
}
