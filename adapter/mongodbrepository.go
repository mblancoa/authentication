package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"time"
)

const (
	CredentialsCollection string = "credentials"
	UserCollection               = "users"
)

type UserDB struct {
	Id          string    `bson:"_id"`
	Email       string    `bson:"email"`
	PhoneNumber string    `bson:"phone_number"`
	Roles       []string  `bson:"roles,omitempty"`
	Last        time.Time `bson:"last,omitempty"`
	Version     int64     `bson:"version,omitempty"`
}

type CredentialsDB struct {
	Id            string               `bson:"_id"`
	Password      string               `bson:"password"`
	State         core.CredentialState `bson:"state"`
	LastPasswords []string             `bson:"last_passwords,omitempty"`
	Attempts      int                  `bson:"attempts"`
	Last          time.Time            `bson:"last"`
	Version       int64                `bson:"version"`
}

//go:generate repogen -dest=mongodbcredentialsrepository_impl.go -model=CredentialsDB -repo=MongoDbCredentialsRepository
type MongoDbCredentialsRepository interface {
	InsertOne(ctx context.Context, credentials *CredentialsDB) (interface{}, error)
	FindById(ctx context.Context, id string) (*CredentialsDB, error)
	UpdateById(ctx context.Context, credentials *CredentialsDB, id string) (bool, error)
}

//go:generate repogen -dest=mongodbuserrepository_impl.go -model=UserDB -repo=MongoDbUserRepository
type MongoDbUserRepository interface {
	FindById(ctx context.Context, id string) (*UserDB, error)
	FindByEmail(ctx context.Context, email string) (*UserDB, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*UserDB, error)
	InsertOne(ctx context.Context, user *UserDB) (interface{}, error)
	UpdateById(ctx context.Context, user *UserDB, id string) (bool, error)
	UpdateEmailById(ctx context.Context, email, id string) (bool, error)
	UpdatePhoneNumberById(ctx context.Context, phoneNumber, id string) (bool, error)
}
