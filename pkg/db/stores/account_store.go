package stores

import (
	"context"
	"fmt"

	"github.com/ochiengotieno304/oneotp/pkg/db"
	"github.com/ochiengotieno304/oneotp/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStore interface {
	CreateAccount(account *models.Account) error
	FindOneAccount(phone string) (*models.Account, error)
	UpdateOneAccount(account *models.Account) error
}

type accountStore struct {
	collection *mongo.Collection
}

func NewAccountStore() AccountStore {
	return &accountStore{
		collection: db.MongoClient().Collection("accounts"),
	}
}

func (store *accountStore) CreateAccount(account *models.Account) error {
	result, err := store.collection.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "name", Value: account.Name},
			{Key: "phone", Value: account.Phone},
			{Key: "email", Value: account.Email},
		},
	)

	fmt.Println(result.InsertedID)

	if err != nil {
		return err
	}

	return nil
}

func (store *accountStore) FindOneAccount(phone string) (*models.Account, error) {
	var account *models.Account
	filter := bson.D{{Key: "phone", Value: phone}}
	err := store.collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (store *accountStore) UpdateOneAccount(account *models.Account) error {
	filter := bson.D{{Key: "phone", Value: account.Phone}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: account.Name},
		{Key: "email", Value: account.Email},
	}}}

	result, err := store.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	fmt.Println(result.UpsertedID)

	return nil
}
