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
	if err := app.authUser(r); err != nil {
		app.errorLog.Println("home: Error authenticating user: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.render(w, r, "home", nil)
}

func (app *application) getApiContent(url string, templateData interface{}) (*http.Response, error) {
	app.infoLog.Printf("Getting content from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	app.infoLog.Printf("Response body: %v/tTemplatedata: %v\n", resp.Body, templateData)
	err = json.NewDecoder(resp.Body).Decode(templateData)
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
	resp, err := http.Post(url, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}

	app.infoLog.Printf("Response body: %v\n", resp.Body)

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

func (app *application) authUser(r *http.Request) error {
	tokenString, err := auth.ExtractToken(r)
	if err != nil {
		return err
	}

	token, err := auth.ParseToken(tokenString)
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	return nil
}