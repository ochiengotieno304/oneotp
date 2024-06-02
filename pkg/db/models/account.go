package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credentials struct {
	SecretKey string `bson:"secret_key,omitempty"`
	APIKey    string `bson:"api_key,omitempty"`
}

type Status struct {
	Code int64 `bson:"code"`
}

type Account struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name,omitempty"`
	Phone        string             `bson:"phone,omitempty"`
	Email        string             `bson:"email,omitempty"`
	Status       Status             `bson:"status"`
	PasswordHash string             `bson:"password_hash"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
	Credentials  []Credentials      `bson:"credentials,omitempty"`
}
