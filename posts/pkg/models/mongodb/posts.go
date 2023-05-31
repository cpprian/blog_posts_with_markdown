package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type PostModel struct {
	C *mongo.Collection
}

// Create a new post
func (m *PostModel) Create(title, content string, userID string) (string, error) {
	return "", nil
}

// Get a post by its ID
func (m *PostModel) GetById(id string) (string, error) {
	return "", nil
}

// Get a post by title
func (m *PostModel) GetByTitle(title string) (string, error) {
	return "", nil
}

// Get all posts
func (m *PostModel) GetAll() (string, error) {
	return "", nil
}

// Get all posts by a user
func (m *PostModel) GetAllByUser(userID string) (string, error) {
	return "", nil
}
