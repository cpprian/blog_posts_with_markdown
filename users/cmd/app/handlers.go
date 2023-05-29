package main

import (
	"encoding/json"
	"net/http"

	"github.com/cpprian/blog_posts_with_markdown/users/pkg/models"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Get all users
	users, err := app.users.All()
	if err != nil {
		app.serverError(w, err)
	}

	// Convert user list into json encoding
	b, err := json.Marshal(users)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("All users were sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findById(w http.ResponseWriter, r *http.Request) {
	// Get user id from request
	id := r.URL.Query().Get("id")

	// Get user
	user, err := app.users.FindById(id)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.serverError(w, err)
	}

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("User was sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByUsername(w http.ResponseWriter, r *http.Request) {
	// Get username from request
	username := r.URL.Query().Get("username")

	// Get user
	user, err := app.users.FindByUsername(username)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.serverError(w, err)
	}

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("User was sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByEmail(w http.ResponseWriter, r *http.Request) {
	// Get email from request
	app.infoLog.Println("Email:", r.URL.Query().Get("email"))
	email := r.URL.Query().Get("email")

	// Get user
	user, err := app.users.FindByEmail(email)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.serverError(w, err)
	}

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("User was sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {
	// Get user from request
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
	}
	app.infoLog.Println("User:", user)

	// Insert user
	_, err = app.users.InsertUser(&user)
	if err != nil {
		app.errorLog.Println("Error:", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app.infoLog.Println("User was inserted with id:", user.ID)

	// Send response
	w.WriteHeader(http.StatusOK)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	// Get user from request
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
	}

	// Update user
	_, err = app.users.UpdateUser(&user)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("User was updated with id:", user.ID)

	// Send response
	w.WriteHeader(http.StatusOK)
}