package repositories

import (
	"context"
	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{collection: collection}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepositoryImpl) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepositoryImpl) PromoteUser(ctx context.Context, username string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"role": "admin"}})
	return err
}

func (r *UserRepositoryImpl) IsFirstUser(ctx context.Context) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	return count == 0, err
}

// In repositories/user_repository.go
func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	cursor, err := r.collection.Find(ctx, bson.M{}) // Get all users
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
