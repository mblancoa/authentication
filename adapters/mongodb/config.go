package mongodb

import (
	"context"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core"
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

func SetupMongodbConfiguration() {
	log.Info().Msg("Initializing mongodb configuration")
	var config mongoDbConfiguration
	core.LoadYamlConfiguration(core.ConfigurationFile, &config)

	conn := config.Mongodb.Database.Connection
	connectionString := fmt.Sprintf("%s:%s//%s:%d", conn.Username, os.Getenv(conn.Password), conn.Host, conn.Port)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	errors.ManageErrorPanic(err)
	err = client.Ping(ctx, nil)
	errors.ManageErrorPanic(err)

	database := client.Database(config.Mongodb.Database.Name)
	setupPersistenceContext(database)
	setupMongoDBMappers()
}

func setupPersistenceContext(database *mongo.Database) {
	mongoDbCredentialsRepository := NewMongoDbCredentialsRepository(database.Collection(CredentialsCollection))
	mongoDbUserRepository := NewMongoDbUserRepository(database.Collection(UserCollection))
	persistenceCtx := core.PersistenceContext
	persistenceCtx.CredentialsPersistenceService = NewMongoDbCredentialsService(mongoDbCredentialsRepository)
	persistenceCtx.UserPersistenceService = NewMongoDbUserService(mongoDbUserRepository)
}

func setupMongoDBMappers() {
	err := mapper.Register(&UserDB{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&CredentialsDB{})
	errors.ManageErrorPanic(err)
}
