package core

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mblancoa/authentication/core/errors"
	"github.com/mblancoa/authentication/core/tools"
	"text/template"
	"time"
)

const (
	// TODO extract Secret to a config file
	Secret    string = "a_*2$ñ6^fjz=?v^66€y7|~2ç"
	SecretJwt string = "$68tk91we&hzDyDhJe[Zz[{&"
)

// TODO review that, is it needed to be an interface?
type AuthenticationService interface {
	// Login checks the authenticity of the user and returns its jws
	Login(credentials UserCredentials) (string, error)

	// ValidateJWT validates a given jwt
	ValidateJWT(jwt string) (bool, error)
	StartSingUP(user User) error
	SingUP(confirmation ConfirmationCredentials) error
}

type authenticationService struct {
	notificationsTemplates *template.Template
	notificationService    tools.NotificationService
	credentialsRepository  UserCredentialsRepository
	userRepository         UserRepository
}

func NewAuthenticationService(notificationService tools.NotificationService, userCredentialsRepository UserCredentialsRepository,
	userRepository UserRepository) AuthenticationService {
	service := authenticationService{
		notificationsTemplates: template.Must(template.ParseGlob("templates/*.txt")),
		notificationService:    notificationService,
		credentialsRepository:  userCredentialsRepository,
		userRepository:         userRepository,
	}
	return &service
}

func (a *authenticationService) Login(credentials UserCredentials) (string, error) {
	var hashedCredentials UserCredentials
	tools.MarshalHash(credentials, &hashedCredentials)
	state, ok := a.credentialsRepository.ExistsUserCredentialsByIdAndPassword(hashedCredentials)
	if !ok {
		return "", errors.NewAuthenticationError("Credentials not found")
	}
	if state.State == Blocked {
		return "", errors.NewAuthenticationError("User Blocked")
	}

	id, err := tools.Encrypt(credentials.ID, Secret)
	if err != nil {
		return "", err
	}

	user, err := a.userRepository.FindUserById(id)
	if err != nil {
		return "", err
	}

	var decUser User
	err = tools.UnmarshalCrypt(user, &decUser, Secret)
	if err != nil {
		return "", err
	}

	return getJwt(decUser, SecretJwt)
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

func (a *authenticationService) StartSingUP(user User) error {
	//TODO implement me
	panic("implement me")
}

func (a *authenticationService) SingUP(confirmation ConfirmationCredentials) error {
	//TODO implement me
	panic("implement me")
}

func getJwt(user User, key string) (string, error) {
	claims := CustomClaims{
		jwt.StandardClaims{
			Id:        user.ID,
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(10)).Unix(),
		},
		user.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
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
