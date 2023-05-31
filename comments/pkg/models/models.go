package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	USER_ID   primitive.ObjectID `bson:"user_id,omitempty"`
	POST_ID   primitive.ObjectID `bson:"post_id,omitempty"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
}
