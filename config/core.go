package config

import (
	"encoding/json"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const (
	RunMode = "RUN_MODE"
)

// ports must be initialized in the configuration of their implementation
var credentialsPersistenceService core.CredentialsPersistenceService
var userPersistenceService core.UserPersistenceService

var notificationService tools.NotificationService
var authenticationService core.AuthenticationService

var configurationFile string = getConfigFile()

func SetupCoreConfiguration() {
	log.Info().Msg("Initializing core configuration")
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

func getConfigFile() string {
	filename := "conf/application.yml"
	if os.Getenv(RunMode) != "" {
		filename = fmt.Sprintf("conf/%s.application.yml", os.Getenv(RunMode))
	}
	return filename
}

func loadJsonConfiguration(fileName string, configObj interface{}) {
	bts := loadFile(fileName)
	err := json.Unmarshal(bts, configObj)
	errors.ManageErrorPanic(err)
}

func loadYamlConfiguration(fileName string, configObj interface{}) {
	bts := loadFile(fileName)
	err := yaml.Unmarshal(bts, configObj)
	errors.ManageErrorPanic(err)
}

func loadFile(file string) []byte {
	confFile, err := os.Open(file)
	errors.ManageErrorPanic(err)
	defer func() {
		err := confFile.Close()
		errors.ManageErrorPanic(err)
	}()

	bts, err := io.ReadAll(confFile)
	errors.ManageErrorPanic(err)
	return bts
}
