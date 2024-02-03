package adapter

import (
	"github.com/mblancoa/authentication/core"
	"time"
)

const (
	CredentialsCollection string = "credentials"
	UserCollection               = "users"
)

type UserDB struct {
	ID          string    `bson:"id"`
	Email       string    `bson:"email"`
	PhoneNumber string    `bson:"phone_number"`
	Roles       []string  `bson:"roles,omitempty"`
	Last        time.Time `bson:"last"`
}

type CredentialsDB struct {
	ID            string               `bson:"id"`
	Password      string               `bson:"password"`
	State         core.CredentialState `bson:"state"`
	LastPasswords []string             `bson:"last_passwords,omitempty"`
	Attempts      int                  `bson:"attempts"`
	Last          time.Time            `bson:"last"`
}
