package adapter

import (
	"github.com/mblancoa/authentication/core"
	"time"
)

type Collection string

const (
	CredentialsCollection Collection = "credentials"
	UserCollection                   = "users"
)

type UserDB struct {
	ID          string    `bson:"id"`
	Email       string    `bson:"email"`
	PhoneNumber string    `bson:"phone_number"`
	Roles       []string  `bson:"roles,omitempty"`
	Last        time.Time `bson:"last"`
}
type Credentials struct {
	ID            string               `bson:"id"`
	Password      string               `bson:"password"`
	State         core.CredentialState `bson:"state"`
	LastPasswords []string             `bson:"last_passwords,omitempty"`
	Attempts      int                  `bson:"attempts"`
	Last          time.Time            `bson:"last"`
}
