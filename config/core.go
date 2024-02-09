package config

import (
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
)

// ports must be initialized in the configuration of their implementation
var credentialsPersistenceService core.CredentialsPersistenceService
var userPersistenceService core.UserPersistenceService

var notificationService tools.NotificationService
var authenticationService core.AuthenticationService

func SetupCoreConfiguration() {
	// TODO initialize notificationService
	authenticationService = core.NewAuthenticationService(notificationService, credentialsPersistenceService, userPersistenceService)

	setupCoreMappers()
}

func setupCoreMappers() {
	err := mapper.Register(&core.FullCredentials{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&core.Credentials{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&core.User{})
	errors.ManageErrorPanic(err)
}
