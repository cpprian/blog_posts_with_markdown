package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home", nil)
}

func (app *application) getApiContent(url string, templateData interface{}) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
