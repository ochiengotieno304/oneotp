package stores

import (
	"context"

	"github.com/ochiengotieno304/oneotp/pkg/db"
	"github.com/ochiengotieno304/oneotp/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthStore interface {
	CreateOTP(otp *models.OTP) (interface{}, error)
	FindOne(id, clientID string) (*models.OTP, error)
	UpdateOne(id, clientID string) error
}

type otpStore struct {
	collection *mongo.Collection
}

func NewOTPStore() AuthStore {
	return &otpStore{
		collection: db.MongoClient().Collection("otp"),
	}
}

func (store *otpStore) CreateOTP(otp *models.OTP) (interface{}, error) {
	result, err := store.collection.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "phone", Value: otp.Phone},
			{Key: "code", Value: otp.Code},
			{Key: "expires_at", Value: otp.ExpiresAt},
			{Key: "used", Value: otp.Used},
			{Key: "client_id", Value: otp.ClientID},
		},
	)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (store *otpStore) FindOne(id, clientID string) (*models.OTP, error) {
	var otp *models.OTP

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{Key: "$and", Value: []interface{}{
			bson.D{{Key: "client_id", Value: clientID}},
			bson.D{{Key: "_id", Value: objID}},
		}},
	}

	err = store.collection.FindOne(
		context.TODO(),
		filter,
	).Decode(&otp)

	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (store *otpStore) UpdateOne(id, clientID string) error { // update only used state
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "$and", Value: []interface{}{
			bson.D{{Key: "client_id", Value: clientID}},
			bson.D{{Key: "_id", Value: objID}},
		}},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "used", Value: true},
		}},
	}

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
