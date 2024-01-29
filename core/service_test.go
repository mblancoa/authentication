package core

// Basic imports
import (
	"github.com/brianvoe/gofakeit"
	"github.com/dgrijalva/jwt-go"
	"github.com/mblancoa/authentication/core/tools"
	"github.com/mblancoa/authentication/errors"
	tools2 "github.com/mblancoa/authentication/tools"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	tools2.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials
	returnedCredentials.State = Active
	var user User
	gofakeit.Struct(&user)
	user.ID = credentials.ID
	user.Roles = []string{"admin", "customer"}
	var encUser User
	_ = tools2.MarshalCrypt(user, &encUser, Secret)

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)
	suite.userRepository.EXPECT().FindUserById(encUser.ID).Return(encUser, nil)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(wToken)

	token, _ := decodeJWT(wToken, SecretJwt)
	suite.Assert().Equal(user.ID, token.Id)
	suite.Assert().Equal(user.Roles, token.Roles)

}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenErrorUnmarshaling() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools2.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Active
	var user User
	gofakeit.Struct(&user)
	id, _ := tools2.Encrypt(credentials.ID, Secret)
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
	tools2.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Active
	id, _ := tools2.Encrypt(credentials.ID, Secret)
	expectedError := errors.NewNotFoundError("User not found")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)
	suite.userRepository.EXPECT().FindUserById(id).Return(User{}, expectedError)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenUserIsBlocked() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools2.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials

	returnedCredentials.State = Blocked
	expectedError := errors.NewAuthenticationError("User Blocked")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(returnedCredentials, true)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenNotExistsUserCredentialsByIdAndPassword() {
	var credentials UserCredentials
	gofakeit.Struct(&credentials)
	var hashCredentials UserCredentials
	tools2.MarshalHash(credentials, &hashCredentials)
	expectedError := errors.NewAuthenticationError("Credentials not found")

	suite.credentialsRepository.EXPECT().ExistsUserCredentialsByIdAndPassword(hashCredentials).Return(UserCredentials{}, false)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_successful() {
	var user User
	gofakeit.Struct(&user)
	j, _ := getJwtFromUser(user, SecretJwt)

	result, err := suite.authenticationService.ValidateJWT(j)

	suite.Assert().NoError(err)
	suite.Assert().True(result)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_failWhenTokenIsNotAnJWT() {
	var wToken string
	gofakeit.Struct(&wToken)

	result, err := suite.authenticationService.ValidateJWT(wToken)

	suite.Assert().Error(err)
	suite.Assert().Equal("token contains an invalid number of segments", err.Error())
	suite.Assert().False(result)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_failWhenTokenIsExpired() {
	id := gofakeit.Username()

	claims := CustomClaims{
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Minute * -1).Unix(),
		},
		[]string{"admin"},
	}
	wToken, _ := getJwt(claims, SecretJwt)

	result, err := suite.authenticationService.ValidateJWT(wToken)

	suite.Assert().Error(err)
	suite.Assert().Regexp("token is expired by .+", err.Error())
	suite.Assert().False(result)
}

func (suite *AuthenticationServiceSuite) TestRefreshJWT_successful() {
	id := gofakeit.Username()

	claims := CustomClaims{
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
		},
		[]string{"admin"},
	}
	wToken, _ := getJwt(claims, SecretJwt)

	rToken, err := suite.authenticationService.RefreshJWT(wToken)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(rToken)
	suite.Assert().NotEqual(wToken, rToken)
}
