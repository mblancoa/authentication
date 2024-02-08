package config

import (
	"encoding/json"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
	"io"
	"os"
)

// ports must be initialized in the configuration of their implementation
var credentialsPersistenceService core.CredentialsPersistenceService
var userPersistenceService core.UserPersistenceService

var notificationService tools.NotificationService
var authenticationService core.AuthenticationService

func SetupCoreConfiguration() {
	// TODO initialize notificationService
	authenticationService = core.NewAuthenticationService(notificationService, credentialsPersistenceService, userPersistenceService)

	SetupCoreMappers()
}

func SetupCoreMappers() {
	err := mapper.Register(&core.FullCredentials{})
	manageErrorPanic(err)
	err = mapper.Register(&core.Credentials{})
	manageErrorPanic(err)
	err = mapper.Register(&core.User{})
	manageErrorPanic(err)
}

func loadConfiguration(configObj any, file string) {
	confFile, err := os.Open(file)
	manageErrorPanic(err)
	defer func() {
		err := confFile.Close()
		manageErrorPanic(err)
	}()

	bts, err := io.ReadAll(confFile)
	manageErrorPanic(err)

	err = json.Unmarshal(bts, configObj)
	manageErrorPanic(err)
}

func manageErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
