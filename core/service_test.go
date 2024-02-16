package core

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mblancoa/authentication/core/domain"
	"github.com/mblancoa/authentication/core/ports"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"github.com/pioz/faker"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type AuthenticationServiceSuite struct {
	suite.Suite
	notificationService           *ports.MockNotificationService
	credentialsPersistenceService *ports.MockCredentialsPersistenceService
	userPersistenceService        *ports.MockUserPersistenceService
	authenticationService         AuthenticationService
}

func (suite *AuthenticationServiceSuite) SetupSuite() {
	_ = os.Chdir("./..")
	suite.notificationService = ports.NewMockNotificationService(suite.T())
	suite.credentialsPersistenceService = ports.NewMockCredentialsPersistenceService(suite.T())
	suite.userPersistenceService = ports.NewMockUserPersistenceService(suite.T())
	suite.authenticationService = NewAuthenticationService(suite.notificationService, suite.credentialsPersistenceService, suite.userPersistenceService)
}

func TestAuthenticationServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationServiceSuite))
}

func (suite *AuthenticationServiceSuite) TestLogin_successful() {
	var credentials domain.Credentials
	tools.FakerBuild(&credentials)
	var hashCredentials domain.Credentials
	tools.MarshalHash(credentials, &hashCredentials)

	returnedCredentials := hashCredentials
	returnedCredentials.State = domain.Active

	var user domain.User
	tools.FakerBuild(&user)
	user.Id = credentials.Id
	user.Roles = []string{"admin", "customer"}
	var encUser domain.User
	_ = tools.MarshalCrypt(user, &encUser, Secret)

	savedCdt := domain.NewFullCredentials(hashCredentials)
	savedCdt.State = domain.Active
	savedCdt.Attempts = 1

	cdtToSave := savedCdt
	cdtToSave.Attempts = 0

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashCredentials.Id).Return(savedCdt, nil)
	suite.credentialsPersistenceService.EXPECT().UpdateCredentials(cdtToSave).Return(nil)
	suite.userPersistenceService.EXPECT().FindUserById(encUser.Id).Return(encUser, nil)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(wToken)

	token, _ := decodeJWT(wToken, SecretJwt)
	suite.Assert().Equal(user.Id, token.Id)
	suite.Assert().Equal(user.Roles, token.Roles)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenErrorUnmarshaling() {
	var credentials domain.Credentials
	tools.FakerBuild(&credentials)
	var hashCredentials domain.Credentials
	tools.MarshalHash(credentials, &hashCredentials)

	savedCdt := domain.NewFullCredentials(hashCredentials)
	savedCdt.State = domain.Active
	savedCdt.Attempts = 0

	var user domain.User
	tools.FakerBuild(&user)
	id, _ := tools.Encrypt(credentials.Id, Secret)
	user.Id = id

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashCredentials.Id).Return(savedCdt, nil)
	suite.userPersistenceService.EXPECT().FindUserById(id).Return(user, nil)

	u, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal("Error decrypting field Email", err.Error())
	suite.Assert().Empty(u)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenCredentialsTurnsBlocked() {
	credentials := domain.Credentials{}
	tools.FakerBuild(&credentials)
	hashed := domain.Credentials{}
	tools.MarshalHash(credentials, &hashed)

	saved := hashed
	saved.Password = faker.UUID()
	saved.Attempts = 2
	saved.State = domain.Active

	toUpdate := domain.FullCredentials{Id: saved.Id, Password: saved.Password, Attempts: 3, State: domain.Blocked}
	expectedError := errors.NewAuthenticationError("Credentials not found")

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashed.Id).Return(domain.NewFullCredentials(saved), nil)
	suite.credentialsPersistenceService.EXPECT().UpdateCredentials(toUpdate).Return(nil)

	result, err := suite.authenticationService.Login(credentials)

	suite.Assertions.Error(err)
	suite.Assertions.Empty(result)
	suite.Assertions.Equal(expectedError, err)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenUserIsBlocked() {
	credentials := domain.Credentials{}
	tools.FakerBuild(&credentials)
	hashed := domain.Credentials{}
	tools.MarshalHash(credentials, &hashed)

	saved := hashed
	saved.State = domain.Blocked

	expectedError := errors.NewAuthenticationError("User Blocked")
	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashed.Id).Return(domain.NewFullCredentials(saved), nil)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}
func (suite *AuthenticationServiceSuite) TestLogin_failWhenCredentialsNotFoundInDB() {
	credentials := domain.Credentials{}
	tools.FakerBuild(&credentials)
	hashed := domain.Credentials{}
	tools.MarshalHash(credentials, &hashed)

	returnedError := tools.NewTestError("Credentials not Found")
	expectedError := errors.NewAuthenticationError("Credentials not found")
	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashed.Id).Return(domain.FullCredentials{}, returnedError)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenUpdatingStateFails() {
	credentials := domain.Credentials{}
	tools.FakerBuild(&credentials)
	hashed := domain.Credentials{}
	tools.MarshalHash(credentials, &hashed)

	saved := hashed
	saved.Password = faker.UUID()
	saved.Attempts = 2
	saved.State = domain.Active

	toUpdate := domain.FullCredentials{Id: hashed.Id, Password: saved.Password, Attempts: 3, State: domain.Blocked}
	returnedError := tools.NewTestError("Error updating")
	expectedError := errors.NewGenericErrorByCause("Error updating credentials attempts", returnedError)

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashed.Id).Return(domain.NewFullCredentials(saved), nil)
	suite.credentialsPersistenceService.EXPECT().UpdateCredentials(toUpdate).Return(returnedError)

	result, err := suite.authenticationService.Login(credentials)

	suite.Assertions.Error(err)
	suite.Assertions.Empty(result)
	suite.Assertions.Equal(expectedError, err)
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenUpdatingAtremptsFails() {
	var credentials domain.Credentials
	tools.FakerBuild(&credentials)
	var hashCredentials domain.Credentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials
	returnedCredentials.State = domain.Active
	var user domain.User
	tools.FakerBuild(&user)
	user.Id = credentials.Id
	user.Roles = []string{"admin", "customer"}
	var encUser domain.User
	_ = tools.MarshalCrypt(user, &encUser, Secret)

	savedCdt := domain.NewFullCredentials(hashCredentials)
	savedCdt.Password = faker.UUID()
	savedCdt.State = domain.Active
	savedCdt.Attempts = 1

	cdtToSave := savedCdt
	cdtToSave.Attempts = 2

	returnedError := tools.NewTestError("Error updating")
	expectedError := errors.NewGenericErrorByCause("Error updating credentials attempts", returnedError)

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashCredentials.Id).Return(savedCdt, nil)
	suite.credentialsPersistenceService.EXPECT().UpdateCredentials(cdtToSave).Return(returnedError)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assertions.Error(err)
	suite.Assertions.Empty(wToken)
	suite.Assertions.Equal(expectedError.Error(), err.Error())
}

func (suite *AuthenticationServiceSuite) TestLogin_failWhenFindUserByIdReturnsError() {
	var credentials domain.Credentials
	tools.FakerBuild(&credentials)
	var hashCredentials domain.Credentials
	tools.MarshalHash(credentials, &hashCredentials)
	returnedCredentials := hashCredentials
	returnedCredentials.State = domain.Active

	savedCdt := domain.NewFullCredentials(hashCredentials)
	savedCdt.State = domain.Active
	savedCdt.Attempts = 1

	cdtToSave := savedCdt
	cdtToSave.Attempts = 0

	userId, _ := tools.Encrypt(credentials.Id, Secret)

	returnedError := errors.NewNotFoundError("User not found")
	expectedError := errors.NewGenericErrorByCause(fmt.Sprintf("Error finding user %s", userId), returnedError)

	suite.credentialsPersistenceService.EXPECT().FindCredentialsById(hashCredentials.Id).Return(savedCdt, nil)
	suite.credentialsPersistenceService.EXPECT().UpdateCredentials(cdtToSave).Return(nil)
	suite.userPersistenceService.EXPECT().FindUserById(userId).Return(domain.User{}, returnedError)

	wToken, err := suite.authenticationService.Login(credentials)

	suite.Assert().Error(err)
	suite.Assert().Equal(expectedError, err)
	suite.Assert().Empty(wToken)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_successful() {
	var user domain.User
	tools.FakerBuild(&user)
	j, _ := getJwtFromUser(user, SecretJwt)

	result, err := suite.authenticationService.ValidateJWT(j)

	suite.Assert().NoError(err)
	suite.Assert().True(result)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_failWhenTokenIsNotAnJWT() {
	var wToken string
	tools.FakerBuild(&wToken)

	result, err := suite.authenticationService.ValidateJWT(wToken)

	suite.Assert().Error(err)
	suite.Assert().Equal("token contains an invalid number of segments", err.Error())
	suite.Assert().False(result)
}

func (suite *AuthenticationServiceSuite) TestValidateJWT_failWhenTokenIsExpired() {
	id := faker.Username()

	claims := domain.CustomClaims{
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
	id := faker.Username()

	claims := domain.CustomClaims{
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
