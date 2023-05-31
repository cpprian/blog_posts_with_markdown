package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/cpprian/blog_posts_with_markdown/posts/pkg/models"
	"github.com/gorilla/mux"
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
	app.render(w, r, "post/editor", nil)
}

func (app *application) createPostPost(w http.ResponseWriter, r *http.Request) {
	var p models.Post
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	p.Title = r.PostForm.Get("title")
	p.Content = string(mdToHTML([]byte(r.PostForm.Get("markdown"))))
	p.CreatedAt = string(time.Now().Format("2006-01-02 15:04:05"))

	app.infoLog.Println("New post ", p)
	if err = app.postApiContent(app.apis.posts, p); err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.getAllPosts(w, r)
}

func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		app.errorLog.Println("Error getting post id")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf(("Getting post with id %s\n"), id)
	url := fmt.Sprintf("%s/%s", app.apis.posts, id)
	app.getPost(w, r, url)
}

func (app *application) getPostByTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title, ok := vars["title"]
	if !ok {
		app.errorLog.Println("Error getting post title")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Printf(("Getting post with title %s\n"), title)
	url := fmt.Sprintf("%s/%s", app.apis.posts, title)
	app.getPost(w, r, url)
}

func (app *application) getPost(w http.ResponseWriter, r *http.Request, url string) {
	app.verifyCookie(w, r)

	resp, err := app.getApiContent(url)
	if err != nil {
		app.errorLog.Println("Error getting post: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var ptd postTempalteData
	err = json.NewDecoder(resp.Body).Decode(&ptd.Post)
	if err != nil {
		app.errorLog.Println("Error decoding post: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.render(w, r, "posts/post", ptd)
}

func (app *application) getAllPosts(w http.ResponseWriter, r *http.Request) {
	app.verifyCookie(w, r)

	app.infoLog.Println("URL: ", app.apis.posts)
	resp, err := app.getApiContent(app.apis.posts)
	if err != nil {
		app.errorLog.Println("Error getting posts: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var ptd postTempalteData
	err = json.NewDecoder(resp.Body).Decode(&ptd.Posts)
	if err != nil {
		app.errorLog.Println("Error decoding posts: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// https://onlinetool.io/goplayground/#txO7hJ-ibeU
func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}