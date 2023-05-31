package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()
	
	// home
	router.HandleFunc("/", app.home).Methods("GET")

	// user handler
	router.HandleFunc("/register", app.registerUser).Methods("GET", "POST")
	router.HandleFunc("/login", app.loginUser).Methods("GET", "POST")
	router.HandleFunc("/logout", app.logoutUser).Methods("GET")
	router.HandleFunc("/user/view/{id:[0-9]+}", app.getUserById).Methods("GET")
	router.HandleFunc("/user/view/{username}", app.getUserByUsername).Methods("GET")

	// post handler
	router.HandleFunc("/posts", app.getAllPosts).Methods("GET")
	router.HandleFunc("/post/create", app.createPost).Methods("GET", "POST")
	router.HandleFunc("/post/view/{id:[0-9]+}", app.getPostById).Methods("GET")
	router.HandleFunc("/post/view/{title}", app.getPostByTitle).Methods("GET")

	// comment handler
	router.HandleFunc("/comment/create", app.createComment).Methods("POST")
	router.HandleFunc("/comment/view/{id:[0-9]+}", app.getCommentsByPostId).Methods("GET")

	router.PathPrefix("/static/").Handler(app.static("./ui/static/"))
	return router
}