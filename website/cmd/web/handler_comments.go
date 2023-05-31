package main

import (
	"net/http"

	"github.com/cpprian/blog_posts_with_markdown/comments/pkg/models"
)

type commentTemplateData struct {
	Comment models.Comment
	Comments []models.Comment
}

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getCommentsByPostId(w http.ResponseWriter, r *http.Request) {

}