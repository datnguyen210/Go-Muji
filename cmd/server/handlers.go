package main

import (
	"errors"
	"fmt"
	"text/template"

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
	}

	for _, blog := range blogs {
		fmt.Fprintf(w, "%+v\n", blog)
	}
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

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Blog: blog,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
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
