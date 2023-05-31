package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	USER_ID primitive.ObjectID `bson:"user_id,omitempty"`
	Title string `bson:"title,omitempty"`
	Content string `bson:"content,omitempty"`
	CreatedAt string `bson:"created_at,omitempty"`
}