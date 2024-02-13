package controllers

import (
	"encoding/json"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
	"net/http"
)

type AuthenticationController struct {
	BaseController
	AuthenticationController core.AuthenticationService
}

func (c *AuthenticationController) Login() {
	var rq LoginRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rq)

	if err != nil {
		c.manageStatusFromError(c.Ctx.ResponseWriter, err)
		return
	}

	credentials := core.Credentials{}
	tools.Mapper(&rq, &credentials)

	authenticationService := WebApplicationContext.AuthenticationService

	jwt, err := authenticationService.Login(credentials)
	if err != nil {
		c.manageStatusFromError(c.Ctx.ResponseWriter, err)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	err = c.Ctx.JSONResp(LoginResponse{Token: jwt})
	if err != nil {
		c.manageStatusFromError(c.Ctx.ResponseWriter, err)
	}
}
