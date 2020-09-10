package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// user := app.extractUser(r)

	if app.runtimeEnv == "dev" {
		fmt.Fprintf(w, "It works! [dev]")
	} else {
		fmt.Fprintf(w, "It works!")
	}
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	countryCode := r.PostForm.Get("country_code")
	number := r.PostForm.Get("number")

	err = app.api.SignUp(countryCode, number, app.clockworkAPI, app.runtimeEnv)
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "1")
}

func (app *application) verifyHash(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	countryCode := r.PostForm.Get("country_code")
	number := r.PostForm.Get("number")
	hash := r.PostForm.Get("hash")

	err = app.api.VerifyHash(countryCode, number, hash)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "1")
}

func (app *application) verifyPin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	countryCode := r.PostForm.Get("country_code")
	number := r.PostForm.Get("number")
	pin := r.PostForm.Get("pin")

	sha1, err := app.api.VerifyPin(countryCode, number, pin)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, sha1)
}
