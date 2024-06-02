package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name,omitempty"`
	Phone     string             `bson:"phone,omitempty"`
	Email     string             `bson:"email,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}
