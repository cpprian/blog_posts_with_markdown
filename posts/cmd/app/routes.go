package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/posts/", app.createPost).Methods("POST")
	router.HandleFunc("/api/posts/{id}", app.getPostById).Methods("GET")
	router.HandleFunc("/api/posts/title/{title}", app.getPostByTitle).Methods("GET")
	router.HandleFunc("/api/posts/", app.getAllPosts).Methods("GET")
	router.HandleFunc("/api/posts/user/{userID}", app.getAllPostsByUser).Methods("GET")

	return router
}