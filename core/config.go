package core

import (
	"encoding/json"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/core/domain"
	"github.com/mblancoa/authentication/core/ports"
	"github.com/mblancoa/authentication/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const (
	RunMode = "RUN_MODE"
)

var ApplicationContext *Context = &Context{}
var configFile string

type Context struct {
	AuthenticationService AuthenticationService
}

func SetupCoreConfig() {
	log.Info().Msg("Initializing core configuration")
	setupCoreContext()
	setupCoreMappers()
}

func setupCoreContext() {
	p := ports.PersistenceContext
	n := ports.NotificationContext
	ApplicationContext.AuthenticationService = NewAuthenticationService(n.NotificationService, p.CredentialsPersistenceService, p.UserPersistenceService)
}

func setupCoreMappers() {
	err := mapper.Register(&domain.FullCredentials{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&domain.Credentials{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&domain.User{})
	errors.ManageErrorPanic(err)
}

func GetConfigFile() string {
	if configFile == "" {
		mode := os.Getenv(RunMode)
		if mode != "" {
			configFile = fmt.Sprintf("conf/%s.application.yml", mode)
		} else {
			configFile = "conf/application.yml"
		}
	}
	return configFile
}

func LoadJsonConfiguration(fileName string, configObj interface{}) {
	bts := loadFile(fileName)
	err := json.Unmarshal(bts, configObj)
	errors.ManageErrorPanic(err)
}

func LoadYamlConfiguration(fileName string, configObj interface{}) {
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
