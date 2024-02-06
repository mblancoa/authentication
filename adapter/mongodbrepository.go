package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CredentialsCollection string = "credentials"
	UserCollection               = "users"
)

type UserDB struct {
	Id          primitive.ObjectID `bson:"_id"`
	UserId      string             `bson:"user_id"`
	Email       string             `bson:"email"`
	PhoneNumber string             `bson:"phone_number"`
	Roles       []string           `bson:"roles,omitempty"`
	Last        time.Time          `bson:"last"`
	Version     int64              `bson:"version"`
}

type CredentialsDB struct {
	Id            primitive.ObjectID   `bson:"_id"`
	UserId        string               `bson:"user_id"`
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
	FindById(ctx context.Context, id primitive.ObjectID) (*CredentialsDB, error)
	FindByUserId(ctx context.Context, userId string) (*CredentialsDB, error)
	UpdateByUserId(ctx context.Context, credentials *CredentialsDB, userId string) (bool, error)
}

//go:generate repogen -dest=mongodbuserrepository_impl.go -model=UserDB -repo=MongoDbUserRepository
type MongoDbUserRepository interface {
	FindById(ctx context.Context, id primitive.ObjectID) (*UserDB, error)
	FindByUserId(ctx context.Context, userId string) (*UserDB, error)
	FindByEmail(ctx context.Context, email string) (*UserDB, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*UserDB, error)
	InsertOne(ctx context.Context, user *UserDB) (interface{}, error)
	UpdateById(ctx context.Context, user *UserDB, id primitive.ObjectID) (bool, error)
	UpdateEmailById(ctx context.Context, email string, id primitive.ObjectID) (bool, error)
	UpdatePhoneNumberById(ctx context.Context, phoneNumber string, id primitive.ObjectID) (bool, error)
}
