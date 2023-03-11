package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/pkg/forms"
	"snippetbox/pkg/models"
	"strconv"
	"strings"
)

var flash = "flash"
var authenticatedUserID = "authenticatedUserID"
var MaxLength = 100
var MinLength = 10
var MaxEmailLength = 254

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl.html", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	app.render(w, r, "show.page.tmpl.html", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var titleName = "title"
	var contentName = "content"
	var expiresName = "expires"

	form := forms.New(r.PostForm)
	form.Values.Set(titleName, strings.TrimSpace(form.Values.Get(titleName)))
	form.Values.Set(contentName, strings.TrimSpace(form.Values.Get(contentName)))
	form.Values.Set(expiresName, strings.TrimSpace(form.Values.Get(expiresName)))
	form.Required(titleName, contentName, expiresName)
	form.MaxLength(titleName, MaxLength)
	form.PermittedValues(expiresName, "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl.html", &templateData{Form: form})
		return
	}
	id, err := app.snippets.Insert(form.Values.Get(titleName), form.Values.Get(contentName), form.Values.Get(expiresName))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, flash, "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var errorEmail = "Address is already in use"
	var successfulSignup = "Your signup was successful. Please log in."

	var nameName = "name"
	var emailName = "email"
	var passwordName = "password"
	form := forms.New(r.PostForm)
	form.Values.Set(nameName, strings.TrimSpace(form.Values.Get(nameName)))
	form.Values.Set(emailName, strings.TrimSpace(form.Values.Get(emailName)))
	form.Required(nameName, emailName, passwordName)
	form.MaxLength(nameName, MaxLength)
	form.MaxLength(emailName, MaxEmailLength)
	form.MatchesPattern(emailName)
	form.MinLength(passwordName, MinLength)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl.html", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Values.Get(nameName), form.Values.Get(emailName), form.Values.Get(passwordName))

	if err != nil {

		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add(emailName, errorEmail)
			app.render(w, r, "signup.page.tmpl.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, flash, successfulSignup)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)

	id, err := app.users.Authenticate(form.Values.Get("email"), form.Values.Get("password"))

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl.html", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, authenticatedUserID, id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, authenticatedUserID)
	app.session.Put(r, flash, "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
