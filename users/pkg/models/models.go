package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the database
type User struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Username   string               `bson:"username"`
	Email      string               `bson:"email"`
	Password   string               `bson:"password"`
	Name       string               `bson:"name"`
	Surname    string               `bson:"surname"`
	About      string               `bson:"about"`
	Posts      []primitive.ObjectID `bson:"posts"`
	Comments   []primitive.ObjectID `bson:"comments"`
	Subscribes []primitive.ObjectID `bson:"subscribe"`
}
