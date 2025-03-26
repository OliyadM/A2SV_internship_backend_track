package data

import (
	"context"
	"errors"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.User) error {
	if userCollection == nil {
		return errors.New("database not initialized")
	}

	
	if existing, _ := FindUserByUsername(user.Username); existing != nil {
		return errors.New("username already exists")
	}

	
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// Set role (first user = admin)
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	_, err = userCollection.InsertOne(context.TODO(), user)
	return err
}

func FindUserByUsername(username string) (*models.User, error) {
	if userCollection == nil {
		return nil, errors.New("database not initialized")
	}

	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func PromoteUser(username string) error {
	if userCollection == nil {
		return errors.New("database not initialized")
	}

	_, err := userCollection.UpdateOne(
		context.TODO(),
		bson.M{"username": username},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	return err
}