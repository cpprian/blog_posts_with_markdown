package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the database
type User struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Username   string               `bson:"username"`
	Email      string               `bson:"email"`
	Password   string               `bson:"password"`
	Subscribes []primitive.ObjectID `bson:"subscribe"`
}
