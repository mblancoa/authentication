package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbCredentialsRepository interface {
	ExistsCredentialsByIdAndPassword(ctx context.Context, credentials CredentialsDB) (CredentialsDB, bool)
	InsertCredentials(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error)
	FindCredentialsById(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error)
	UpdateCredentials(ctx context.Context, credentials CredentialsDB) error
}

type mongoDbCredentialsRepository struct {
	collection *mongo.Collection
}

func NewMongoDbCredentialsRepository(database *mongo.Database) MongoDbCredentialsRepository {
	collection := database.Collection(CredentialsCollection)
	return &mongoDbCredentialsRepository{collection: collection}
}

func (m *mongoDbCredentialsRepository) ExistsCredentialsByIdAndPassword(ctx context.Context, credentials CredentialsDB) (CredentialsDB, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbCredentialsRepository) InsertCredentials(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbCredentialsRepository) FindCredentialsById(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbCredentialsRepository) UpdateCredentials(ctx context.Context, credentials CredentialsDB) error {
	//TODO implement me
	panic("implement me")
}

type MongoDbUserRepository interface {
	FindUserById(ctx context.Context, id string) (UserDB, error)
	FindUserByEmail(ctx context.Context, email string) (UserDB, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (UserDB, error)
	InsertUser(ctx context.Context, user core.User) (UserDB, error)
	UpdateUser(ctx context.Context, user core.User) (UserDB, error)
}

type mongoDbUserRepository struct {
	collection *mongo.Collection
}

func NewMongoDbUserRepository(database *mongo.Database) MongoDbUserRepository {
	collection := database.Collection(UserCollection)
	return &mongoDbUserRepository{collection: collection}
}

func (m *mongoDbUserRepository) FindUserById(ctx context.Context, id string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbUserRepository) FindUserByEmail(ctx context.Context, email string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbUserRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbUserRepository) InsertUser(ctx context.Context, user core.User) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDbUserRepository) UpdateUser(ctx context.Context, user core.User) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}
