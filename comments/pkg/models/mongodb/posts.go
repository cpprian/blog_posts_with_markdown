package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type CommentModel struct {
	C *mongo.Collection
}
