package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	// Create a new mux router
	router := mux.NewRouter()

	// Register handlers
	router.HandleFunc("/api/users/", app.all).Methods("GET")
	router.HandleFunc("/api/users/{id}", app.findById).Methods("GET")
	router.HandleFunc("/api/users/{username}", app.findByUsername).Methods("GET")
	router.HandleFunc("/api/users/email/{email}", app.findByEmail).Methods("GET")
	router.HandleFunc("/api/users/", app.insertUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", app.updateUser).Methods("PUT")

	return router
}