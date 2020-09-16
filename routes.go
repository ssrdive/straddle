package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	r := mux.NewRouter()
	r.Handle("/", http.HandlerFunc(app.home)).Methods("GET")
	r.Handle("/signup", http.HandlerFunc(app.signUp)).Methods("POST")
	r.Handle("/verifyPin", http.HandlerFunc(app.verifyPin)).Methods("POST")
	r.Handle("/updateProfile", http.HandlerFunc(app.updateProfile)).Methods("POST")
	r.Handle("/getDetails/{number}", http.HandlerFunc(app.getDetails)).Methods("GET")
	r.Handle("/verifyHash", http.HandlerFunc(app.verifyHash)).Methods("POST")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
