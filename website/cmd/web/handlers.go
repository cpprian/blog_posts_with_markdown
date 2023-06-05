package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cpprian/blog_posts_with_markdown/website/pkg/auth"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Getting home page")
	app.verifyCookie(w, r)
	app.getAllPosts(w, r)
}

func (app *application) getApiContent(url string) (*http.Response, error) {
	app.infoLog.Printf("Getting content from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (app *application) postApiContent(url string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	app.infoLog.Printf("Posting content %v to %s\n", data, url)
	_, err = http.Post(url, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}

	return nil
}

func (app *application) static(dir string) http.Handler {
	dirCleaned := filepath.Clean(dir)
	return http.StripPrefix("/static/", http.FileServer(http.Dir(dirCleaned)))
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td interface{}) {
	files := []string{
		"./ui/html/" + name + ".page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) authUser(tokenString string) error {
	token, err := auth.ParseToken(tokenString)
	if err != nil {
		app.errorLog.Println("authUser: Error parsing token: ", err)
		return err
	}

	app.infoLog.Println("authUser: Token is ", token , " and is ", token.Valid)
	if !token.Valid {
		app.errorLog.Println("authUser: Token is not valid")
		return err
	}

	return nil
}

func (app *application) verifyCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		app.errorLog.Println("verifyCookie: Error getting cookie: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenString := cookie.Value
	if err := app.authUser(tokenString); err != nil {
		app.errorLog.Println("verifyCookie: Error authenticating user: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}