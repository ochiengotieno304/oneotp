package models

import "time"

type Credentials struct {
	SecretKey string `bson:"secret_key,omitempty"`
}

type Status struct {
	Code int64 `bson:"code"`
}

var (
	Unverified Status = Status{
		Code: 100,
	}

	Verified Status = Status{
		Code: 200,
	}
)

type Account struct {
	ID           string      `bson:"_id"`
	Name         string      `bson:"name,omitempty"`
	Phone        string      `bson:"phone,omitempty"`
	Email        string      `bson:"email,omitempty"`
	Status       int64       `bson:"status"`
	PasswordHash string      `bson:"password_hash"`
	CreatedAt    time.Time   `bson:"created_at,omitempty"`
	UpdatedAt    time.Time   `bson:"updated_at,omitempty"`
	Credentials  Credentials `bson:"credentials,omitempty"`
}
