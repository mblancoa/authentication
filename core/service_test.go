package core

// Basic imports
import (
	"github.com/brianvoe/gofakeit"
	"github.com/mblancoa/authentication/core/errors"
	"github.com/mblancoa/authentication/core/tools"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthenticationServiceSuite struct {
	suite.Suite
	notificationService   *tools.MockNotificationService
	credentialsRepository *MockUserCredentialsRepository
	userRepository        *MockUserRepository
	authenticationService AuthenticationService
}

func (suite *AuthenticationServiceSuite) SetupTest() {
	suite.notificationService = tools.NewMockNotificationService(suite.T())
	suite.credentialsRepository = NewMockUserCredentialsRepository(suite.T())
	suite.userRepository = NewMockUserRepository(suite.T())
	suite.authenticationService = NewAuthenticationService(suite.notificationService, suite.credentialsRepository, suite.userRepository)
}

func TestAuthenticationServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationServiceSuite))
}

func (suite *AuthenticationServiceSuite) TestLogin_successful() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials
	returnedCredentials.State = Active
	var user User
	gofakeit.Struct(&user)
	user.ID = credentials.ID
	var encUser User
	_ = tools.MarshalCrypt(user, &encUser, Secret)

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)
	suite.userRepository.EXPECT().FindUserById(encUser.ID).Return(encUser, nil)

	jwt, err := suite.authenticationService.Login(credentials)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(jwt)

	token, _ := decodeJWT(jwt, SecretJwt)
	suite.Assert().Equal(user.ID, token.Id)
	suite.Assert().Equal(user.Roles, token.Roles)

}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenErrorUnmarshaling() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Active
	var user User
	gofakeit.Struct(&user)
	id, _ := tools.Encrypt(credentials.ID, Secret)
	user.ID = id
	expectedError := "Error decrypting field Email"

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)
	suite.userRepository.EXPECT().FindUserById(id).Return(user, nil)

	u, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err.Error())
	suite.Assert().Empty(u)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenFindUserByIdReturnsError() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Active
	id, _ := tools.Encrypt(credentials.ID, Secret)
	expectedError := errors.NewNotFoundError("User not found")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)
	suite.userRepository.EXPECT().FindUserById(id).Return(User{}, expectedError)

	jwt, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(jwt)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenUserIsBlocked() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Blocked
	expectedError := errors.NewAuthenticationError("User Blocked")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)

	jwt, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(jwt)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenNotExistsUserCredentialsByIdAndPassword() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools.MarshalHash(credentials, &hashCredentials)
	expectedError := errors.NewAuthenticationError("Credentials not found")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(UserCredentials{}, false)

	jwt, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(jwt)
}
