package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/datnguyen210/go-blog/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.blogs.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Blogs = blogs

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) viewBlog(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(r)
	data.Blog = blog

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) modalCreateBlog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) createBlog(w http.ResponseWriter, r *http.Request) {
	// Parse form data from the request
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Retrieve values from form fields
	title := r.FormValue("title")
	content := r.FormValue("content")
	expiresStr := r.FormValue("expires")

	// Convert expires to an integer
	expires, err := strconv.Atoi(expiresStr)
	if err != nil || expires < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Insert the blog into the database
	id, err := app.blogs.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the new blog post's view page
	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
