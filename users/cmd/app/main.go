package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cpprian/blog_posts_with_markdown/users/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    *mongodb.UserModel
}

func main() {
	serverAddr := flag.String("serverAddr", "", "HTTP server network address")
	serverPort := flag.Int("serverPort", 4000, "HTTP server network port")
	mongoURI := flag.String("mongoURI", "mongodb://localhost:27017", "MongoDB connection URI")
	mongoDB := flag.String("mongoDB", "users", "MongoDB database name")
	enableCredentials := flag.Bool("enableCredentials", false, "Enable credentials for mongodb connection")
	flag.Parse()

	infoLog := log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create mongo client
	co := options.Client().ApplyURI(*mongoURI)
	if *enableCredentials {
		co.Auth = &options.Credential{
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		}
	}

	// Establish connection
	client, err := mongo.NewClient(co)
	if err != nil {
		errorLog.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			errorLog.Fatal(err)
		}
	}()

	infoLog.Printf("Connected to MongoDB on %s\n", *mongoURI)

	// Create a new application
	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
		users: &mongodb.UserModel{
			C: client.Database(*mongoDB).Collection("users"),
		},
	}

	// Start the server
	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:         serverURI,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", serverURI)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
