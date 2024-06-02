package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID        primitive.ObjectID `bson:"_id"`
	Phone     string             `bson:"phone,omitempty"`
	Code      string             `bson:"code,omitempty"`
	Used      bool               `bson:"used,omitempty"`
	ExpiresAt time.Time          `bson:"expires_at,omitempty"`
}
