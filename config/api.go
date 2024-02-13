package config

import (
	"github.com/devfeel/mapper"
	"github.com/mblancoa/authentication/api/controllers"
	"github.com/mblancoa/authentication/errors"
	"github.com/rs/zerolog/log"
)

func SetupApiConfiguration() {
	log.Info().Msg("Initializing api configuration")
	controllers.WebApplicationContext.AuthenticationService = authenticationService
	setupApiMappers()
}
func setupApiMappers() {
	err := mapper.Register(&controllers.LoginRequest{})
	errors.ManageErrorPanic(err)
	err = mapper.Register(&controllers.LoginResponse{})
	errors.ManageErrorPanic(err)
}
