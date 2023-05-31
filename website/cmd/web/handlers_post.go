package main

import (
	"net/http"

	"github.com/cpprian/blog_posts_with_markdown/posts/pkg/models"
)

type postTempalteData struct {
	Post models.Post
	Posts []models.Post
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	app.verifyCookie(w, r)

	switch r.Method {
	case "GET":
		app.createPostGet(w, r)
	case "POST":
		app.createPostPost(w, r)
	default:
		app.errorLog.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *application) createPostGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "posts/create", nil)
}

func (app *application) createPostPost(w http.ResponseWriter, r *http.Request) {
	// var p models.Post
	// err := r.ParseForm()
	// if err != nil {
	// 	app.errorLog.Println(err.Error())
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

}

func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getPostByTitle(w http.ResponseWriter, r *http.Request) {

}