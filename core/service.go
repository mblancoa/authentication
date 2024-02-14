package core

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"text/template"
	"time"
)

const (
	// TODO extract to a config file
	Secret      string = "a_*2$ñ6^fjz=?v^66€y7|~2ç"
	SecretJwt   string = "$68tk91we&hzDyDhJe[Zz[{&"
	MaxAttempts int    = 3
)

// TODO review that, is it needed to be an interface?
type AuthenticationService interface {
	// Login checks the authenticity of the user and returns its jws
	Login(credentials Credentials) (string, error)

	// ValidateJWT validates a given jwt
	ValidateJWT(jwt string) (bool, error)

	// RefreshJWT generates a new jwt
	RefreshJWT(jwt string) (string, error)

	StartSingUP(user User) error
	SingUP(confirmation ConfirmationCredentials) error
}

type authenticationService struct {
	notificationsTemplates        *template.Template
	notificationService           tools.NotificationService
	credentialsPersistenceService CredentialsPersistenceService
	userPersistenceService        UserPersistenceService
}

func NewAuthenticationService(notificationService tools.NotificationService, credentialsPersistenceService CredentialsPersistenceService,
	userPersistenceService UserPersistenceService) AuthenticationService {
	service := authenticationService{
		notificationsTemplates:        template.Must(template.ParseGlob("./templates/*.txt")),
		notificationService:           notificationService,
		credentialsPersistenceService: credentialsPersistenceService,
		userPersistenceService:        userPersistenceService,
	}
	return &service
}

func (a *authenticationService) Login(credentials Credentials) (string, error) {
	var hashedCredentials Credentials
	tools.MarshalHash(credentials, &hashedCredentials)

	state, err := a.checkCredentials(hashedCredentials, MaxAttempts)
	if err != nil {
		if errors.GetCode(err) == errors.NotFoundError {
			return "", errors.NewAuthenticationError("Credentials not found")
		}
		return "", err
	}
	if state.State == Blocked {
		return "", errors.NewAuthenticationError("User Blocked")
	}

	userId, err := tools.Encrypt(credentials.Id, Secret)
	if err != nil {
		return "", err
	}

	user, err := a.userPersistenceService.FindUserById(userId)
	if err != nil {
		return "", errors.NewGenericErrorByCause(fmt.Sprintf("Error finding user %s", userId), err)
	}

	var decUser User
	err = tools.UnmarshalCrypt(user, &decUser, Secret)
	if err != nil {
		return "", err
	}

	return getJwtFromUser(decUser, SecretJwt)
}

func (a *authenticationService) ValidateJWT(token string) (bool, error) {
	tokenData, err := decodeJWT(token, SecretJwt)
	if err != nil {
		return false, err
	}

	err = tokenData.Valid()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (a *authenticationService) RefreshJWT(token string) (string, error) {
	claims, err := decodeJWT(token, SecretJwt)
	if err != nil {
		return "", err
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(10)).Unix()
	refresh, err := getJwt(*claims, SecretJwt)
	if err != nil {
		return "", errors.NewGenericError("Error refreshing token")
	}
	return refresh, nil

}

func (a *authenticationService) StartSingUP(user User) error {
	//TODO implement me
	panic("implement me")
}

func (a *authenticationService) SingUP(confirmation ConfirmationCredentials) error {
	//TODO implement me
	panic("implement me")
}

// checkCredentials finds credentials by userId and checks its password.
// If credentials are not found an error is returned
// If password is incorrect attempts get increased and an error is returned. If attempts are equal to max
// credentials get blocked.
//
// credentials must be hashed
func (a *authenticationService) checkCredentials(credentials Credentials, maxAttempts int) (Credentials, error) {
	savedCredentials, err := a.credentialsPersistenceService.FindCredentialsById(credentials.Id)
	if err != nil {
		return Credentials{}, errors.NewNotFoundError(err.Error())
	}
	if savedCredentials.Password != credentials.Password {
		savedCredentials.Attempts++
		if savedCredentials.Attempts == 3 {
			savedCredentials.State = Blocked
		}
		err = a.credentialsPersistenceService.UpdateCredentials(savedCredentials)
		if err != nil {
			return Credentials{}, errors.NewGenericErrorByCause("Error updating credentials attempts", err)
		}
		return Credentials{}, errors.NewNotFoundError("credentials not found")
	} else if savedCredentials.State == Active && savedCredentials.Attempts != 0 {
		savedCredentials.Attempts = 0
		err = a.credentialsPersistenceService.UpdateCredentials(savedCredentials)
		if err != nil {
			return Credentials{}, errors.NewGenericErrorByCause("Error updating credentials state", err)
		}
	}
	var result Credentials
	tools.Mapper(&savedCredentials, &result)

	return result, nil
}

func getJwt(claims CustomClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func getJwtFromUser(user User, key string) (string, error) {
	claims := CustomClaims{
		jwt.StandardClaims{
			Id:        user.Id,
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(10)).Unix(),
		},
		user.Roles,
	}
	return getJwt(claims, key)
}

func decodeJWT(token string, key string) (*CustomClaims, error) {
	var claims CustomClaims // custom claims
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil // using a config struct to handle the secret
	})

	if err != nil {
		return nil, err
	}

	return &claims, nil
}
