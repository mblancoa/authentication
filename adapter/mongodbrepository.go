package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbCredentialsRepository struct {
	collection *mongo.Collection
}

func NewMongoDbCredentialsRepository(database *mongo.Database) *MongoDbCredentialsRepository {
	collection := database.Collection(CredentialsCollection)
	return &MongoDbCredentialsRepository{collection: collection}
}

func (m *MongoDbCredentialsRepository) ExistsCredentialsByIdAndPassword(ctx context.Context, credentials CredentialsDB) (CredentialsDB, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbCredentialsRepository) InsertCredentials(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbCredentialsRepository) FindCredentialsById(ctx context.Context, credentials CredentialsDB) (CredentialsDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbCredentialsRepository) UpdateCredentials(ctx context.Context, credentials CredentialsDB) error {
	//TODO implement me
	panic("implement me")
}

type MongoDbUserRepository struct {
	collection *mongo.Collection
}

func NewMongoDbUserRepository(database *mongo.Database) *MongoDbUserRepository {
	collection := database.Collection(CredentialsCollection)
	return &MongoDbUserRepository{collection: collection}
}

func (m *MongoDbUserRepository) FindUserById(ctx context.Context, id string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) FindUserByEmail(ctx context.Context, email string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) InsertUser(ctx context.Context, user core.User) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDbUserRepository) UpdateUser(ctx context.Context, user core.User) (UserDB, error) {
	//TODO implement me
	panic("implement me")
}
