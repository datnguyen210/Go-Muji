package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) createBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.blogs.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fieldsErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldsErrors["title"] = "Title cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldsErrors["title"] = "Title exceeds 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		fieldsErrors["content"] = "Content cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldsErrors["expires"] = "Invalid expiration date: expires must be 1, 7 or 365"
	}

	if len(fieldsErrors) > 0 {
		fmt.Fprintf(w, "Validation errors: %v", fieldsErrors)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}
