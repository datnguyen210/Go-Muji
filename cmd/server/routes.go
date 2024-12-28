package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Automatically load and save session data with every HTTP request and response
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/user/register", dynamic.ThenFunc(app.userRegister))
	router.Handler(http.MethodPost, "/user/register", dynamic.ThenFunc(app.userRegisterPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/blog/view/:id", dynamic.ThenFunc(app.viewBlog))
	router.Handler(http.MethodGet, "/blog/create", dynamic.ThenFunc(app.modalCreateBlog))
	router.Handler(http.MethodPost, "/blog/create", dynamic.ThenFunc(app.createBlog))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
