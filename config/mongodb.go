package config

import (
	"context"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/adapter"
	"github.com/mblancoa/authentication/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type mongoDbConfiguration struct {
	Mongodb struct {
		Database struct {
			Name       string `yaml:"name"`
			Connection struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"connection"`
		} `yaml:"database"`
	} `yaml:"mongodb"`
}

var mongoDbCredentialsRepository adapter.MongoDbCredentialsRepository
var mongoDbUserRepository adapter.MongoDbUserRepository

func SetupMongodbConfiguration() {
	log.Info().Msg("Initializing mongodb configuration")
	var config mongoDbConfiguration
	loadYamlConfiguration(configurationFile, &config)

	conn := config.Mongodb.Database.Connection
	connectionString := fmt.Sprintf("%s:%s//%s:%d", conn.Username, os.Getenv(conn.Password), conn.Host, conn.Port)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	errors.ManageErrorPanic(err)
	err = client.Ping(ctx, nil)
	errors.ManageErrorPanic(err)

	database := client.Database(config.Mongodb.Database.Name)
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
