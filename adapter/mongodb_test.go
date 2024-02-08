package adapter

import (
	"context"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"github.com/pioz/faker"
	"github.com/stretchr/testify/suite"
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
	tools.FakerBuild(&db)
	db.Attempts = 1
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: db.Password}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(db.Id, result.Id)
	suite.Assert().Equal(db.Password, result.Password)
	suite.Assert().Equal(core.Active, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assert().Equal(0, updated.Attempts)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsBlockedCredentiaslWhenCredentialsGetsBlocked() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	db.Attempts = 2
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(core.Blocked, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assert().Equal(3, updated.Attempts)
	suite.Assert().Equal(core.Blocked, updated.State)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsBlockedCredentialsWhenCredentialsWasBlocked() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	db.Attempts = 1
	db.State = core.Blocked

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(core.Blocked, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assert().Equal(2, updated.Attempts)
	suite.Assert().Equal(core.Blocked, updated.State)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsErrorWhenCredentialsNotFound() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assert().Error(err)
	suite.Assert().Equal("mongo: no documents in result", err.Error())
	suite.Assert().Empty(result)
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestInsertCredentials_successful() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)

	result, err := suite.credentialsPersistenceService.InsertCredentials(credentials)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(credentials, result)
	suite.Assert().Equal(int64(1), count(suite.credentialsCollection, context.TODO()))
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestInsertCredentials_returnsErrorWhenCredentialsWithSameUserIdExists() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)
	db := CredentialsDB{}
	_ = mapper.Mapper(&credentials, &db)
	insertOne(suite.credentialsCollection, context.TODO(), &db)

	result, err := suite.credentialsPersistenceService.InsertCredentials(credentials)

	suite.Assert().Error(err)
	suite.Assert().Empty(result)
	suite.Assert().Equal(errors.Error, errors.GetCode(err, ""))
}
