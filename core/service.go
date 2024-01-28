package core

import (
	"github.com/mblancoa/authentication/core/errors"
	"github.com/mblancoa/authentication/core/tools"
	"text/template"
)

// TODO extract Secret to a config file
const Secret string = "a_*2$ñ6^fjz=?v^66€y7|~2ç"

// TODO review that, is it needed to be an interface?
type AuthenticationService interface {
	//TODO return jwt
	Login(credentials UserCredentials) (User, error)
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

func (a *authenticationService) Login(credentials UserCredentials) (User, error) {
	var hashedCredentials UserCredentials
	tools.MarshalHash(credentials, &hashedCredentials)
	state, ok := a.credentialsRepository.ExistsUserCredentialsByIdAndPassword(hashedCredentials)
	if !ok {
		return User{}, errors.NewAuthenticationError("Credentials not found")
	}
	if state.State == Blocked {
		return User{}, errors.NewAuthenticationError("User Blocked")
	}

	id, err := tools.Encrypt(credentials.ID, Secret)
	if err != nil {
		return User{}, err
	}

	user, err := a.userRepository.FindUserById(id)
	if err != nil {
		return User{}, err
	}

	var decUser User
	err = tools.UnMarshalCrypt(user, &decUser, Secret)
	if err != nil {
		return User{}, err
	}

	return decUser, nil
}

func (a *authenticationService) StartSingUP(user User) error {
	//TODO implement me
	panic("implement me")
}

func (a *authenticationService) SingUP(confirmation ConfirmationCredentials) error {
	//TODO implement me
	panic("implement me")
}
