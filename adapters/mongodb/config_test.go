package mongodb

import (
	"context"
	mim "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/core/ports"
	"github.com/mblancoa/authentication/tools"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var mongodbServer *mim.Server

func init() {
	err := os.Chdir("./../..")
	tools.ManageTestError(err)
	err = os.Setenv(core.RunMode, "test")
	tools.ManageTestError(err)
}

func setupDB() {
	server, err := mim.StartWithOptions(context.TODO(), "5.0.2", mim.WithPort(27017))
	tools.ManageTestError(err)
	mongodbServer = server
}

func TearDownDB() {
	mongodbServer.Stop(context.TODO())
}

func TestLoadConfiguration(t *testing.T) {
	var config mongoDbConfiguration
	core.LoadYamlConfiguration(core.GetConfigFile(), &config)

	assert.NotEmpty(t, config)
	assert.NotEmpty(t, config.Mongodb)
	db := config.Mongodb.Database
	assert.NotEmpty(t, db)
	assert.Equal(t, "auth", db.Name)
	con := db.Connection
	assert.NotEmpty(t, con)
	assert.Equal(t, "localhost", con.Host)
	assert.Equal(t, int(27017), con.Port)
	assert.Equal(t, "mongodb", con.Username)
	assert.Equal(t, "TEST_DB_PASSWORD", con.Password)
}

func TestSetupMongodbConfiguration(t *testing.T) {
	setupDB()
	defer TearDownDB()

	SetupMongodbConfiguration()

	assert.NotEmpty(t, ports.PersistenceContext)
	assert.NotEmpty(t, ports.PersistenceContext.UserPersistenceService)
	assert.NotEmpty(t, ports.PersistenceContext.CredentialsPersistenceService)
}
