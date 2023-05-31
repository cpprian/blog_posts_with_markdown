package mongodb

import (
	"context"
	"errors"

	"github.com/cpprian/blog_posts_with_markdown/posts/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostModel struct {
	C *mongo.Collection
}

// Create a new post
func (m *PostModel) Create(post *models.Post) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()

	res, err := m.C.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Get a post by its ID
func (m *PostModel) GetById(id string) (*models.Post, error) {
	ctx := context.TODO()
	var post models.Post
	err := m.C.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &post, nil
}

// Get a post by title
func (m *PostModel) GetByTitle(title string) (*models.Post, error) {
	ctx := context.TODO()
	var post models.Post

	err := m.C.FindOne(ctx, bson.M{"title": title}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &post, nil
}

// Get all posts
func (m *PostModel) GetAll() ([]models.Post, error) {
	ctx := context.TODO()
	var posts []models.Post

	cursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Get all posts by a user
func (m *PostModel) GetAllByUser(userID string) ([]models.Post, error) {
	u, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	var posts []models.Post

	cursor, err := m.C.Find(ctx, bson.M{"user_id": u})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
