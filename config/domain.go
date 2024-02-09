package config

import (
	"encoding/json"
	"fmt"
	"github.com/mblancoa/authentication/errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

const (
	RunMode = "RUN_MODE"
)

type AppConfiguration struct {
	MongoDB MongoDbConfiguration `json:"mongodb" yaml:"mongodb"`
}

// MongoDb configuration domain

type MongoDbConfiguration struct {
	Database DB `json:"database" yaml:"database"`
}
type DB struct {
	Name       string     `json:"name"`
	Connection Connection `json:"connection"`
}
type Connection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var appConfiguration *AppConfiguration = nil

func GetAppConfiguration() *AppConfiguration {
	if appConfiguration == nil {
		filename := "conf/application.yml"
		if os.Getenv(RunMode) != "" {
			filename = fmt.Sprintf("conf/%s.application.yml", os.Getenv(RunMode))
		}
		appConfiguration = &AppConfiguration{}
		loadYamlConfiguration(appConfiguration, filename)
	}
	return appConfiguration
}

func loadJsonConfiguration(configObj any, file string) {
	bts, err := loadFile(file)

	err = json.Unmarshal(bts, configObj)
	errors.ManageErrorPanic(err)
}

func loadYamlConfiguration(configObj any, file string) {
	bts, err := loadFile(file)

	err = yaml.Unmarshal(bts, configObj)
	errors.ManageErrorPanic(err)
}

func loadFile(file string) ([]byte, error) {
	confFile, err := os.Open(file)
	errors.ManageErrorPanic(err)
	defer func() {
		err := confFile.Close()
		errors.ManageErrorPanic(err)
	}()

	bts, err := io.ReadAll(confFile)
	errors.ManageErrorPanic(err)
	return bts, err
}
