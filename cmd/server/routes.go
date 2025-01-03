package main

import (
	"net/http"

	"github.com/datnguyen210/go-blog/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Automatically load and save session data with every HTTP request and response
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/user/register", dynamic.ThenFunc(app.userRegister))
	router.Handler(http.MethodPost, "/user/register", dynamic.ThenFunc(app.userRegisterPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodGet, "/blog/view/:id", dynamic.ThenFunc(app.viewBlog))

	protected := dynamic.Append(app.mustBeAuthenticated)

	router.Handler(http.MethodGet, "/blog/create", protected.ThenFunc(app.modalCreateBlog))
	router.Handler(http.MethodPost, "/blog/create", protected.ThenFunc(app.createBlog))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
