package mongodb

import (
	"context"
	"errors"

	"github.com/cpprian/blog_posts_with_markdown/users/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel wraps a mongo collection and provides methods to query it
type UserModel struct {
	C *mongo.Collection
}

// All returns all users
func (m *UserModel) All() ([]models.User, error) {
	ctx := context.TODO()
	var users []models.User

	cursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID finds a user by id
func (m *UserModel) FindById(id string) (*models.User, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	var user models.User

	if err := m.C.FindOne(ctx, bson.M{"_id": p}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &user, nil
}

// FindByUsername finds a user by username
func (m *UserModel) FindByUsername(username string) (*models.User, error) {
	ctx := context.TODO()
	var user models.User

	if err := m.C.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (m *UserModel) FindByEmail(email string) (*models.User, error) {
	ctx := context.TODO()
	var user models.User

	if err := m.C.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no user found")
		}
		return nil, err
	}

	return &user, nil
}

// InsertUser inserts a new user to the database
func (m *UserModel) InsertUser(user *models.User) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()

	res, err := m.C.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateUser updates the user with the given id when posts, comments or subscribes are added/removed/updated
func (m *UserModel) UpdateUser(user *models.User) (*mongo.UpdateResult, error) {
	ctx := context.TODO()

	res, err := m.C.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		return nil, err
	}

	return res, nil
}
