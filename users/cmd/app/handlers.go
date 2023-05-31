package main

import (
	"encoding/json"
	"net/http"

	"github.com/cpprian/blog_posts_with_markdown/users/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// Get all users
	users, err := app.users.All()
	if err != nil {
		app.errorLog.Println("Error getting all users: ", err)
		app.serverError(w, err)
		return
	}

	// Convert user list into json encoding
	b, err := json.Marshal(users)
	if err != nil {
		app.errorLog.Println("Error marshalling users: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("\nAll users were sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findById(w http.ResponseWriter, r *http.Request) {
	// Get user id from request
	id := mux.Vars(r)["id"]
	app.infoLog.Printf("Getting user with id %s\n", id)

	// Get user
	user, err := app.users.FindById(id)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.errorLog.Println("Error getting user: ", err)
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("\nUser:", user)

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.errorLog.Println("Error marshalling user: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("\nUser was sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByUsername(w http.ResponseWriter, r *http.Request) {
	// Get username from request
	username := mux.Vars(r)["username"]
	app.infoLog.Printf("Getting user with username %s\n", username)

	// Get user
	user, err := app.users.FindByUsername(username)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.errorLog.Println("Error getting user: ", err)
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("\nUser:", user)

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.errorLog.Println("Error marshalling user: ", err)
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("\nUser was sent")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByEmail(w http.ResponseWriter, r *http.Request) {
	// Get email from request
	email := mux.Vars(r)["email"]
	app.infoLog.Println("Email:", email)

	// Get user
	user, err := app.users.FindByEmail(email)
	if err != nil {
		if err.Error() == "no user found" {
			app.infoLog.Println("User not found")
			return
		}
		app.errorLog.Println("Error getting user: ", err)
		app.serverError(w, err)
		return
	}
	app.infoLog.Println("\nUser:", user)

	// Convert user into json encoding
	b, err := json.Marshal(user)
	if err != nil {
		app.errorLog.Println("Error marshalling user: ", err)
		app.serverError(w, err)
		return
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
		app.errorLog.Println("Error:", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.infoLog.Println("\nUser:", user)

	// Insert user
	_, err = app.users.InsertUser(&user)
	if err != nil {
		app.errorLog.Println("Error:", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app.infoLog.Println("User was inserted with data:", user)

	// Send response
	w.WriteHeader(http.StatusOK)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	// Get user from request
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.errorLog.Println("Error:", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.infoLog.Println("\nUser:", user)

	// Update user
	_, err = app.users.UpdateUser(&user)
	if err != nil {
		app.errorLog.Println("Error:", err)
		app.serverError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app.infoLog.Println("User was updated with id:", user.ID)

	// Send response
	w.WriteHeader(http.StatusOK)
}