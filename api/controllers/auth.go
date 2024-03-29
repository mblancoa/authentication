package controllers

import (
	"encoding/json"
	"github.com/mblancoa/authentication/core/domain"
	"github.com/mblancoa/authentication/tools"
	"net/http"
)

type AuthenticationController struct {
	BaseController
}

func (c *AuthenticationController) Login() {
	var rq LoginRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rq)

	if err != nil {
		c.manageStatusFromError(c.Ctx.ResponseWriter, err)
		return
	}

	credentials := domain.Credentials{}
	tools.Mapper(&rq, &credentials)

	jwt, err := c.AuthenticationService.Login(credentials)
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
