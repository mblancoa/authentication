package adapter

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/pioz/faker"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type MongoDBCredentialsPersistenceServiceSuite struct {
	MongoDBPersistenceSuite
	credentialsPersistenceService core.CredentialsPersistenceService
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) SetupSuite() {
	suite.MongoDBPersistenceSuite.SetupSuite()
	suite.setupCredentialsCollection()
	suite.credentialsPersistenceService = NewMongoDbCredentialsService(suite.credentialsRepository)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) SetupTest() {
	ctx := context.Background()
	deleteAll(suite.credentialsCollection, ctx)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TearDownSuite() {
	suite.MongoDBPersistenceSuite.TearDownSuite()
}

func TestCredentialsServiceSuite(t *testing.T) {
	suite.Run(t, new(MongoDBCredentialsPersistenceServiceSuite))
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_successful() {
	db := CredentialsDB{}
	err := faker.Build(&db)
	manageTestError(err)
	db.Attempts = 1
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{UserId: db.UserId, Password: db.Password}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(db.UserId, result.UserId)
	suite.Assert().Equal(db.Password, result.Password)
	suite.Assert().Equal(core.Active, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "user_id", db.UserId, &updated)
	suite.Assert().Equal(0, updated.Attempts)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsErrorWhenCredentialsGetsBlocked() {
	db := CredentialsDB{}
	err := faker.Build(&db)
	manageTestError(err)
	db.Attempts = 2
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{UserId: db.UserId, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().Error(err)
	suite.Assert().Equal("Authentication error", err.Error())
	suite.Assert().Equal("BasicError", reflect.TypeOf(err).Name())
	suite.Assert().Equal(errors.AuthenticationError, err.(errors.BasicError).Code)
	suite.Assert().Empty(result)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "user_id", db.UserId, &updated)
	suite.Assert().Equal(3, updated.Attempts)
	suite.Assert().Equal(core.Blocked, updated.State)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsErrorWhenCredentialsWasBlocked() {
	db := CredentialsDB{}
	err := faker.Build(&db)
	manageTestError(err)
	db.Attempts = 1
	db.State = core.Blocked

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{UserId: db.UserId, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().Error(err)
	suite.Assert().Equal("Authentication error", err.Error())
	suite.Assert().Equal("BasicError", reflect.TypeOf(err).Name())
	suite.Assert().Equal(errors.AuthenticationError, err.(errors.BasicError).Code)
	suite.Assert().Empty(result)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "user_id", db.UserId, &updated)
	suite.Assert().Equal(2, updated.Attempts)
	suite.Assert().Equal(core.Blocked, updated.State)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsErrorWhenCredentialsNotFound() {
	credentials := core.Credentials{}
	err := faker.Build(&credentials)
	manageTestError(err)

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)
	suite.Assert().Equal("Authentication error\nCaused by mongo: no documents in result", err.Error())
	suite.Assert().Equal("BasicError", reflect.TypeOf(err).Name())
	suite.Assert().Equal(errors.AuthenticationError, err.(errors.BasicError).Code)

	suite.Assert().Error(err)
	suite.Assert().Empty(result)
}
