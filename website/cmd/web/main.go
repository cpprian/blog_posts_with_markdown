package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type apis struct {
	users string
	posts string
	comments string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	apis     apis
}

func main() {
	serverAddr := flag.String("serverAddr", "", "HTTP server network address")
	serverPort := flag.Int("serverPort", 8000, "HTTP server network port")
	usersAPI := flag.String("usersAPI", "http://localhost:4000/api/users", "Users API endpoint")
	postsAPI := flag.String("postsAPI", "http://localhost:4000/api/posts", "Posts API endpoint")
	commentsAPI := flag.String("commentsAPI", "http://localhost:4000/api/comments", "Comments API endpoint")
	flag.Parse()

	app := &application{
		errorLog: log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime),
		apis: apis{
			users: *usersAPI,
			posts: *postsAPI,
			comments: *commentsAPI,
		},
	}

	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:     serverURI,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", serverURI)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}