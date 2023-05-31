package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/posts", app.createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", app.getPostById).Methods("GET")
	router.HandleFunc("/posts/title/{title}", app.getPostByTitle).Methods("GET")
	router.HandleFunc("/posts", app.getAllPosts).Methods("GET")
	router.HandleFunc("/posts/user/{userID}", app.getAllPostsByUser).Methods("GET")

	return router
}