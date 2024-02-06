package adapter

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/mblancoa/authentication/core"
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

func (suite *MongoDBCredentialsPersistenceServiceSuite) TearDownSuite() {
	suite.MongoDBPersistenceSuite.TearDownSuite()
}

func TestCredentialsServiceSuite(t *testing.T) {
	suite.Run(t, new(MongoDBCredentialsPersistenceServiceSuite))
}

func (suite *MongoDBCredentialsPersistenceServiceSuite) TestFindCredentialsById() {
	var credentials CredentialsDB
	gofakeit.Struct(&credentials)
	gofakeit.Struct(&credentials.LastPasswords)
	insertOne(suite.credentialsCollection, context.TODO(), credentials)

	result, err := suite.credentialsRepository.FindByUserId(context.TODO(), credentials.UserId)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(result)
	suite.Assert().Equal(&credentials, result)
}
