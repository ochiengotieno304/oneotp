package models

import "time"

type OTP struct {
	ID        string    `bson:"_id"`
	Phone     string    `bson:"phone,omitempty"`
	Code      string    `bson:"code,omitempty"`
	Used      bool      `bson:"used,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
	ClientID  string    `bson:"client_id,omitempty"`
}
