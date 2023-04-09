package main

import (
	"net/http"

	// New import
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, app.noSurf, app.authenticate)

	mux := pat.New()

	// snippets
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// user
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// ping
	mux.Get("/ping", http.HandlerFunc(ping))

	// static
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}

// func (app *application) routes() http.Handler {
// 	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

// 	fileServer := http.FileServer(http.Dir("./ui/static"))
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", app.home)
// 	mux.HandleFunc("/snippet", app.showSnippet)
// 	mux.HandleFunc("/snippet/create", app.createSnippet)
// 	mux.Handle("/home", &myHome{})
// 	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

// 	return standardMiddleware.Then(mux)
// 	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
// 	// return secureHeaders(mux)
// }
