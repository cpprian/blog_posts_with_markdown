package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cpprian/blog_posts_with_markdown/users/pkg/models"
	"github.com/cpprian/blog_posts_with_markdown/website/pkg/auth"
	"github.com/gorilla/mux"
)

type userTemplateData struct {
	User models.User
	Users []models.User
}

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		app.errorLog.Println("Error getting user id")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf(("Getting user with id %s\n"), id)
	url := fmt.Sprintf("%s/%s", app.apis.users, id)

	app.getUser(w, r, url)
}

func (app *application) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		app.errorLog.Println("Error getting user username")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf(("Getting user with username %s\n"), username)
	url := fmt.Sprintf("%s/%s", app.apis.users, username)
	app.getUser(w, r, url)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request, url string) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.errorLog.Println("Error getting token")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var utd userTemplateData
	app.getApiContent(url, &utd.Users)
	app.infoLog.Println(utd.Users)

	app.render(w, r, "users", utd)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.registerUserGet(w, r)
	case "POST":
		app.registerUserPost(w, r)
	default:
		app.errorLog.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *application) registerUserGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register", nil)
}

func (app *application) registerUserPost(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	u.Username = strings.TrimSpace(r.PostForm.Get("username"))
	u.Email = strings.TrimSpace(r.PostForm.Get("email"))
	u.Password = r.PostForm.Get("password")

	u.Password, err = auth.EncryptPassword(u.Password)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/", app.apis.users)
	app.postApiContent(url, u)

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.loginUserGet(w, r)
	case "POST":
		app.loginUserPost(w, r)
	default:
		app.errorLog.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *application) loginUserGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login", nil)
}

func (app *application) loginUserPost(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	u.Username = strings.TrimSpace(r.PostForm.Get("email"))
	u.Password = strings.TrimSpace(r.PostForm.Get("password"))

	resp, err := app.getApiContent(fmt.Sprintf("%s/email/%s", app.apis.users, u.Username), &u)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := json.NewDecoder(resp.Body)
	var authUser models.User
	err = user.Decode(&authUser)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = auth.ComparePassword(authUser.Password, u.Password)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := auth.NewToken(authUser.ID.Hex())
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", token)

	app.infoLog.Println("User logged in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.errorLog.Println("Error getting token")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}