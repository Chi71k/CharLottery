package repository

import (
	"context"
	"user-service/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	collection := client.Database("userdb").Collection("users")
	return &UserRepository{collection: collection}
}

func (repo *UserRepository) UserExists(username string) bool {
	filter := bson.M{"username": username}
	count, err := repo.collection.CountDocuments(context.Background(), filter)
	return err == nil && count > 0
}

func (repo *UserRepository) CreateUser(user *model.User) error {
	_, err := repo.collection.InsertOne(context.Background(), user)
	return err
}

func (repo *UserRepository) GetUserByID(userID string) (*model.User, error) {
	filter := bson.M{"_id": userID}
	var user model.User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	filter := bson.M{"username": username}
	var user model.User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	filter := bson.M{"email": email}
	var user model.User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdatePassword(userID, newPassword string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"password": newPassword}}
	_, err := repo.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (repo *UserRepository) GetAllUsers() ([]model.User, error) {
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []model.User
	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(userID, username, email string) error {
    filter := bson.M{"_id": userID}
    update := bson.M{"$set": bson.M{"username": username, "email": email}}
    _, err := r.collection.UpdateOne(context.Background(), filter, update)
    return err
}

func (r *UserRepository) DeleteUser(userID string) error {
    filter := bson.M{"_id": userID}
    _, err := r.collection.DeleteOne(context.Background(), filter)
    return err
}