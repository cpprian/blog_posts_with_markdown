package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()
	// home
	router.HandleFunc("/", app.home).Methods("GET")
	// register 
	router.HandleFunc("/register", app.registerUser).Methods("GET", "POST")
	// login
	router.HandleFunc("/login", app.loginUser).Methods("GET", "POST")
	// logout
	router.HandleFunc("/logout", app.logoutUser).Methods("GET")
	// get user
	router.HandleFunc("/user/view/{id:[0-9]+}", app.getUserById).Methods("GET")
	router.HandleFunc("/user/view/{username}", app.getUserByUsername).Methods("GET")


	router.PathPrefix("/static/").Handler(app.static("./ui/static/"))
	return router
}