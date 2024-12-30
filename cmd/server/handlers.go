package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/datnguyen210/go-blog/internal/models"
	"github.com/datnguyen210/go-blog/internal/validator"
	"github.com/julienschmidt/httprouter"
)

var registryTemplate = "registry.tmpl"

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

type blogCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
	// the tag `form:"-"` indicates that this field should be ignored
	// during form processing
}

func (app *application) modalCreateBlog(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = blogCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) createBlog(w http.ResponseWriter, r *http.Request) {
	var form blogCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title exceeds 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "Invalid expiration date: expires must be 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.blogs.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Blog created successfully ✅")

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}

type userResgistryForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
	// the tag `form:"-"` indicates that this field should be ignored
	// during form processing
}

func (app *application) userRegister(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userResgistryForm{}
	app.render(w, http.StatusOK, registryTemplate, data)
}

func (app *application) userRegisterPost(w http.ResponseWriter, r *http.Request) {
	var form userResgistryForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "Username cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please enter a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, registryTemplate, data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplidatedEmail) {
			form.AddfieldError("email", "This email already exists")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, registryTemplate, data)
		} else {
			app.serverError(w, err)
		}

		app.sessionManager.Put(r.Context(), "flash", "Register successfully, please log in")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Please enter a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is not correct")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserId")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
