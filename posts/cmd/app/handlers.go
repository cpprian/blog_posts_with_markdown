package main

import (
	"encoding/json"
	"net/http"

	"github.com/cpprian/blog_posts_with_markdown/posts/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Create a new post")

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		app.errorLog.Println("Error decoding JSON: ", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app.infoLog.Println(post)

	_, err = app.posts.Create(&post)
	if err != nil {
		app.errorLog.Println("Error creating post: ", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.infoLog.Println("Post created successfully")
	w.WriteHeader(http.StatusCreated)
}

func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Get a post by its ID")

	id := mux.Vars(r)["id"]
	app.infoLog.Printf("ID: %s", id)

	post, err := app.posts.GetById(id)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("Post not found")
			return
		}
		app.errorLog.Println("Error getting post: ", err)
		app.serverError(w, err)
		return
	}
	app.infoLog.Println(post)

	b, err := json.Marshal(post)
	if err != nil {
		app.errorLog.Println("Error marshalling post: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Post marshalled successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) getPostByTitle(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Get a post by title")

	title := mux.Vars(r)["title"]
	app.infoLog.Printf("Title: %s", title)

	post, err := app.posts.GetByTitle(title)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("Post not found")
			return
		}
		app.errorLog.Println("Error getting post: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(post)

	b, err := json.Marshal(post)
	if err != nil {
		app.errorLog.Println("Error marshalling post: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Post marshalled successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) getAllPosts(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Get all posts")

	posts, err := app.posts.GetAll()
	if err != nil {
		app.errorLog.Println("Error getting posts: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(posts)

	b, err := json.Marshal(posts)
	if err != nil {
		app.errorLog.Println("Error marshalling posts: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Posts marshalled successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) getAllPostsByUser(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Get all posts by a user")

	userID := mux.Vars(r)["userID"]
	app.infoLog.Printf("User ID: %s", userID)

	posts, err := app.posts.GetAllByUser(userID)
	if err != nil {
		app.errorLog.Println("Error getting posts: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(posts)

	b, err := json.Marshal(posts)

	if err != nil {
		app.errorLog.Println("Error marshalling posts: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Posts marshalled successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}