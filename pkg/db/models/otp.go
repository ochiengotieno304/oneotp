package models

import "time"

type OTP struct {
	ID        string    `bson:"_id"`
	Phone     string    `bson:"phone,omitempty"`
	Code      string    `bson:"code,omitempty"`
	Attempts  int8      `bson:"attempts,omitempty"`
	Expired   bool      `bson:"expired,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
	ClientID  string    `bson:"client_id,omitempty"`
}
