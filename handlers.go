package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func (app *application) getDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]
	if number == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	cds, err := app.api.Details(number)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cds)
}

func (app *application) updateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	number := r.PostForm.Get("number")
	firstName := r.PostForm.Get("first_name")
	lastName := r.PostForm.Get("last_name")
	displayName := r.PostForm.Get("display_name")
	dob := r.PostForm.Get("dob")
	status := r.PostForm.Get("status")

	err = app.api.UpdateProfile(number, firstName, lastName, displayName, dob, status)
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
