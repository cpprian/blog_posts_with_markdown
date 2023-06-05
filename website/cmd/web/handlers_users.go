package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	app.verifyCookie(w, r)

	resp, err := app.getApiContent(url)
	if err != nil {
		app.errorLog.Println("Error getting user: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var utd userTemplateData
	err = json.NewDecoder(resp.Body).Decode(&utd.User)
	if err != nil {
		app.errorLog.Println("Error decoding users: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.getAllPosts(w, r)
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
	app.render(w, r, "auth/register", nil)
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

	if err = app.postApiContent(app.apis.users, u); err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	app.render(w, r, "auth/login", nil)
}

func (app *application) loginUserPost(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	u.Email = strings.TrimSpace(r.PostForm.Get("email"))
	u.Password = strings.TrimSpace(r.PostForm.Get("password"))

	resp, err := app.getApiContent(fmt.Sprintf("%semail/%s", app.apis.users, u.Email))
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userDecoder := json.NewDecoder(resp.Body)
	var authUser models.User
	err = userDecoder.Decode(&authUser)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	app.infoLog.Println(authUser)

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

	cookie := http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	app.infoLog.Println("User logged in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name: "token",
		Value: "",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) getAllUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := app.getApiContent(app.apis.users)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var utd userTemplateData
	err = json.NewDecoder(resp.Body).Decode(&utd.Users)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Println(utd.Users)
	app.getAllPosts(w, r)
}

func (app *application) getUsernameById(id string) (string, error) {
	resp, err := app.getApiContent(fmt.Sprintf("%s%s", app.apis.users, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var u models.User
	err = json.NewDecoder(resp.Body).Decode(&u)
	if err != nil {
		return "", err
	}

	return u.Username, nil
}