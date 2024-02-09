package config

import (
	"context"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/adapter"
	"github.com/mblancoa/authentication/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var mongoDbCredentialsRepository adapter.MongoDbCredentialsRepository
var mongoDbUserRepository adapter.MongoDbUserRepository

// SetupMongoDBConfiguration sets mongodb configuration
func SetupMongoDBConfiguration() {
	config := GetAppConfiguration().MongoDB
	conn := config.Database.Connection
	connectionString := fmt.Sprintf("%s:%s//%s:%d", conn.Username, os.Getenv(conn.Password), conn.Host, conn.Port)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	errors.ManageErrorPanic(err)

	database := client.Database(config.Database.Name)
	mongoDbCredentialsRepository = adapter.NewMongoDbCredentialsRepository(database.Collection(adapter.CredentialsCollection))
	mongoDbUserRepository = adapter.NewMongoDbUserRepository(database.Collection(adapter.UserCollection))
	credentialsPersistenceService = adapter.NewMongoDbCredentialsService(mongoDbCredentialsRepository)
	userPersistenceService = adapter.NewMongoDbUserService(mongoDbUserRepository)

	setupMongoDBMappers()
}

func setupMongoDBMappers() {
	err := mapper.Register(&adapter.UserDB{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&adapter.CredentialsDB{})
	errors.ManageErrorPanic(err)
}
