package stores

import (
	"context"
	"time"

	"github.com/ochiengotieno304/oneotp/pkg/db"
	"github.com/ochiengotieno304/oneotp/pkg/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStore interface {
	CreateAccount(account *models.Account) (interface{}, error)
	FindAccountByEmail(email string) error
	FindAccountByID(id string) (*models.Account, error)
	UpdateAccountCredentials(account *models.Account) error
	DeleteOneAccount(id string) error
}

type accountStore struct {
	collection *mongo.Collection
}

func NewAccountStore() AccountStore {
	return &accountStore{
		collection: db.MongoClient().Collection("accounts"),
	}
}

func (store *accountStore) CreateAccount(account *models.Account) (interface{}, error) {
	result, err := store.collection.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "name", Value: account.Name},
			{Key: "phone", Value: account.Phone},
			{Key: "email", Value: account.Email},
			{Key: "password_hash", Value: account.PasswordHash},
			{Key: "created_at", Value: time.Now()},
			{Key: "updated_at", Value: time.Now()},
			{Key: "status", Value: 100},
		},
	)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (store *accountStore) FindAccountByID(id string) (*models.Account, error) {
	var account *models.Account

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	err = store.collection.FindOne(
		context.TODO(),
		filter,
	).Decode(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (store *accountStore) FindAccountByEmail(email string) error {
	filter := bson.D{{Key: "email", Value: email}}
	if err := store.collection.FindOne(context.TODO(), filter).Err(); err != nil {
		return err
	}
	
	return nil
}

func (store *accountStore) UpdateAccountCredentials(account *models.Account) error {
	objID, err := primitive.ObjectIDFromHex(account.ID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "credentials", Value: bson.D{{
			Key: "secret_key", Value: account.Credentials.SecretKey,
		}}},
	}}}

	_, err = store.collection.UpdateOne(
		context.TODO(),
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *accountStore) DeleteOneAccount(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	_, err = store.collection.DeleteOne(
		context.TODO(),
		filter,
	)

	if err != nil {
		return err
	}

	return nil
}
