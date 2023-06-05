package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/cpprian/blog_posts_with_markdown/posts/pkg/models"
	"github.com/cpprian/blog_posts_with_markdown/website/pkg/auth"
	"github.com/gorilla/mux"
)

type PostData struct {
	Post models.Post
	Username string
	Content template.HTML
}

type PostTempalteData struct {
	Post PostData
	Posts []PostData
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
	p.Content = r.PostForm.Get("markdown")
	app.infoLog.Println("Content: ", p.Content)
	p.CreatedAt = string(time.Now().Format("2006-01-02 15:04:05"))

	cookie, err := r.Cookie("token")
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user_id, err := auth.GetUserIdFromToken(cookie.Value)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	p.USER_ID, err = primitive.ObjectIDFromHex(user_id)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.infoLog.Println("New post ", p)
	if err = app.postApiContent(app.apis.posts, p); err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) getPostById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app.infoLog.Println(vars)
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

	var ptd PostTempalteData
	err = json.NewDecoder(resp.Body).Decode(&ptd.Post.Post)
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

	var ptd PostTempalteData
	var posts []models.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		app.errorLog.Println("Error decoding posts: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		username, err := app.getUsernameById(post.USER_ID.Hex())
		if err != nil {
			app.errorLog.Println("Error getting username: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		app.infoLog.Println("Username: ", username)

		content := template.HTML(mdToHTML([]byte(post.Content)))
		app.infoLog.Println("Content: ", content)

		ptd.Posts = append(ptd.Posts, PostData{
			Post: post,
			Username: username,
			Content: content,
		})
	}

	app.render(w, r, "home", ptd)
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