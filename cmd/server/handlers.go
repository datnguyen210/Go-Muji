package main

import (
	"encoding/json"
	"errors"
	"fmt"

	// "html/template"
	"net/http"
	"strconv"

	"github.com/datnguyen210/go-blog/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	blogs, err := app.blogs.Latest()
	if err != nil {
		app.serverError(w, err)
		return 
	}
	for _, blog := range(blogs){
		fmt.Fprintf(w, "%+v\n", blog )
	}


	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...) // destructure
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	blog, err := app.blogs.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode the blog data to JSON and write it to the response
	err = json.NewEncoder(w).Encode(blog)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// dummy data
	title := "Building a second brain"
	content := "Building a second brain with CODE and PARA method"
	expires := 7

	id, err := app.blogs.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect user to the blog details after creation
	http.Redirect(w, r, fmt.Sprintf("/blog/view?id=%d", id), 200)
}
