package controllers

import "github.com/mblancoa/authentication/core"

type webApplicationContext struct {
	AuthenticationService core.AuthenticationService
}

var WebApplicationContext *webApplicationContext = &webApplicationContext{}
