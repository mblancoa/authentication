package core

// Basic imports
import (
	"github.com/brianvoe/gofakeit"
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

	u, err := suite.authenticationService.Login(credentials)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(u)
	suite.Assert().Equal(user, u)
}
