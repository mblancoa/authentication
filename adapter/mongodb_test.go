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

// CredentialsPersistenceService Tests

type mongoDBCredentialsPersistenceServiceSuite struct {
	mongoDBPersistenceSuite
	credentialsPersistenceService core.CredentialsPersistenceService
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) SetupSuite() {
	suite.mongoDBPersistenceSuite.SetupSuite()
	suite.setupCredentialsCollection()
	suite.credentialsPersistenceService = NewMongoDbCredentialsService(suite.credentialsRepository)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) SetupTest() {
	ctx := context.Background()
	deleteAll(suite.credentialsCollection, ctx)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TearDownSuite() {
	suite.mongoDBPersistenceSuite.TearDownSuite()
}

func TestCredentialsServiceSuite(t *testing.T) {
	suite.Run(t, new(mongoDBCredentialsPersistenceServiceSuite))
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_successful() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	db.Attempts = 1
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: db.Password}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(result)
	suite.Assertions.Equal(db.Id, result.Id)
	suite.Assertions.Equal(db.Password, result.Password)
	suite.Assertions.Equal(core.Active, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assertions.Equal(0, updated.Attempts)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsBlockedCredentiaslWhenCredentialsGetsBlocked() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	db.Attempts = 2
	db.State = core.Active

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(result)
	suite.Assertions.Equal(core.Blocked, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assertions.Equal(3, updated.Attempts)
	suite.Assertions.Equal(core.Blocked, updated.State)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsBlockedCredentialsWhenCredentialsWasBlocked() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	db.Attempts = 1
	db.State = core.Blocked

	insertOne(suite.credentialsCollection, context.TODO(), db)
	credentials := core.Credentials{Id: db.Id, Password: faker.String()}

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(result)
	suite.Assertions.Equal(core.Blocked, result.State)

	updated := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", db.Id, &updated)
	suite.Assertions.Equal(2, updated.Attempts)
	suite.Assertions.Equal(core.Blocked, updated.State)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestCheckCredentials_returnsErrorWhenCredentialsNotFound() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)

	result, err := suite.credentialsPersistenceService.CheckCredentials(credentials, 3)

	suite.Assertions.Error(err)
	suite.Assertions.Equal("mongo: no documents in result", err.Error())
	suite.Assertions.Empty(result)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestInsertCredentials_successful() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)

	result, err := suite.credentialsPersistenceService.InsertCredentials(credentials)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(result)
	suite.Assertions.Equal(credentials, result)
	suite.Assertions.Equal(int64(1), count(suite.credentialsCollection, context.TODO()))
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestInsertCredentials_returnsErrorWhenCredentialsWithSameUserIdExists() {
	credentials := core.Credentials{}
	tools.FakerBuild(&credentials)
	db := CredentialsDB{}
	_ = mapper.Mapper(&credentials, &db)
	insertOne(suite.credentialsCollection, context.TODO(), &db)

	result, err := suite.credentialsPersistenceService.InsertCredentials(credentials)

	suite.Assertions.Error(err)
	suite.Assertions.Empty(result)
	suite.Assertions.Equal(errors.Error, errors.GetCode(err))
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestFindCredentialsByUserId_successful() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	insertOne(suite.credentialsCollection, context.TODO(), &db)

	result, err := suite.credentialsPersistenceService.FindCredentialsByUserId(db.Id)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(result)
	credentials := core.FullCredentials{}
	tools.Mapper(&db, &credentials)
	suite.Assertions.Equal(credentials, result)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestFindCredentialsByUserId_returnsErrorWhenNotFound() {
	id := faker.UUID()
	result, err := suite.credentialsPersistenceService.FindCredentialsByUserId(id)

	suite.Assertions.Error(err)
	suite.Assertions.Equal(errors.NotFoundError, errors.GetCode(err))
	suite.Assertions.Empty(result)
}

func (suite *mongoDBCredentialsPersistenceServiceSuite) TestUpdateCredentials_successful() {
	db := CredentialsDB{}
	tools.FakerBuild(&db)
	id := db.Id
	insertOne(suite.credentialsCollection, context.TODO(), &db)

	toUpdate := core.FullCredentials{}
	tools.FakerBuild(&toUpdate)
	toUpdate.Id = id

	err := suite.credentialsPersistenceService.UpdateCredentials(toUpdate)

	suite.Assertions.NoError(err)
	updatedDb := CredentialsDB{}
	findOne(suite.credentialsCollection, context.TODO(), "_id", id, &updatedDb)
	updated := core.FullCredentials{}
	tools.Mapper(&updatedDb, &updated)
	suite.Assertions.Equal(toUpdate, updated)
}

// UserPersistenceService Test

type mongoDBUserPersistenceServiceSuite struct {
	mongoDBPersistenceSuite
	userPersistenceService core.UserPersistenceService
}

func (suite *mongoDBUserPersistenceServiceSuite) SetupSuite() {
	suite.mongoDBPersistenceSuite.SetupSuite()
	suite.setupUserCollection()
	suite.userPersistenceService = NewMongoDbUserService(suite.userRepository)
}

func (suite *mongoDBUserPersistenceServiceSuite) SetupTest() {
	ctx := context.Background()
	deleteAll(suite.userCollection, ctx)
}

func (suite *mongoDBUserPersistenceServiceSuite) TearDownSuite() {
	suite.mongoDBPersistenceSuite.TearDownSuite()
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(mongoDBUserPersistenceServiceSuite))
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserById_successful() {
	db := UserDB{}
	tools.FakerBuild(&db)
	insertOne(suite.userCollection, context.TODO(), &db)

	user, err := suite.userPersistenceService.FindUserById(db.Id)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(user)

	updated := core.User{}
	tools.Mapper(&db, &updated)
	suite.Assertions.Equal(updated, user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserById_returnsErrorWhenNotFound() {
	user, err := suite.userPersistenceService.FindUserById(faker.UUID())

	suite.Assertions.Error(err)
	suite.Assertions.Equal(errors.NotFoundError, errors.GetCode(err))
	suite.Assertions.Equal("mongo: no documents in result", err.Error())
	suite.Assertions.Empty(user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserByEmail_successful() {
	db := UserDB{}
	tools.FakerBuild(&db)
	insertOne(suite.userCollection, context.TODO(), &db)

	user, err := suite.userPersistenceService.FindUserByEmail(db.Email)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(user)

	updated := core.User{}
	tools.Mapper(&db, &updated)
	suite.Assertions.Equal(updated, user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserByEmail_returnsErrorWhenNotFound() {
	user, err := suite.userPersistenceService.FindUserByEmail(faker.Email())

	suite.Assertions.Error(err)
	suite.Assertions.Equal(errors.NotFoundError, errors.GetCode(err))
	suite.Assertions.Equal("mongo: no documents in result", err.Error())
	suite.Assertions.Empty(user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserByPhoneNumber_successful() {
	db := UserDB{}
	tools.FakerBuild(&db)
	insertOne(suite.userCollection, context.TODO(), &db)

	user, err := suite.userPersistenceService.FindUserByPhoneNumber(db.PhoneNumber)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(user)

	updated := core.User{}
	tools.Mapper(&db, &updated)
	suite.Assertions.Equal(updated, user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestFindUserByPhoneNumber_returnsErrorWhenNotFound() {
	user, err := suite.userPersistenceService.FindUserByPhoneNumber(faker.PhoneNumber())

	suite.Assertions.Error(err)
	suite.Assertions.Equal(errors.NotFoundError, errors.GetCode(err))
	suite.Assertions.Equal("mongo: no documents in result", err.Error())
	suite.Assertions.Empty(user)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestInsertUser_successful() {
	user := core.User{}
	tools.FakerBuild(&user)

	returned, err := suite.userPersistenceService.InsertUser(user)

	suite.Assertions.NoError(err)
	suite.Assertions.NotEmpty(returned)
	suite.Assertions.Equal(returned, user)
	suite.Assertions.Equal(int64(1), count(suite.userCollection, context.TODO()))
}

func (suite *mongoDBUserPersistenceServiceSuite) TestInsertUser_returnsErrorWhenUserIdAlreadyExists() {
	db := UserDB{}
	tools.FakerBuild(&db)
	insertOne(suite.userCollection, context.TODO(), &db)

	user := core.User{}
	tools.FakerBuild(&user)
	user.Id = db.Id

	returned, err := suite.userPersistenceService.InsertUser(user)

	suite.Assertions.Error(err)
	suite.Assertions.Equal(errors.Error, errors.GetCode(err))
	suite.Assertions.Empty(returned)
}

func (suite *mongoDBUserPersistenceServiceSuite) TestUpdateUser_successful() {
	db := UserDB{}
	tools.FakerBuild(&db)
	id := db.Id
	insertOne(suite.userCollection, context.TODO(), &db)

	toUpdate := core.User{}
	tools.FakerBuild(&toUpdate)
	toUpdate.Id = id

	err := suite.userPersistenceService.UpdateUser(toUpdate)

	suite.Assertions.NoError(err)
	updatedDb := UserDB{}
	findOne(suite.userCollection, context.TODO(), "_id", id, &updatedDb)
	updated := core.User{}
	tools.Mapper(&updatedDb, &updated)
	suite.Assertions.Equal(toUpdate, updated)
}
