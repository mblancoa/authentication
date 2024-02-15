package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type webApplicationContext struct {
	AuthenticationService core.AuthenticationService
}

type BaseController struct {
	beego.Controller
	*core.Context
}

func (c *BaseController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Context = core.ApplicationContext
}

func (c *BaseController) manageStatusFromError(response http.ResponseWriter, err error) {
	log.Info().Msgf("Managing the error\n%s", err.Error())
	var status int
	switch code := errors.GetCode(err); code {
	case errors.AuthenticationError:
		status = http.StatusUnauthorized
		break
	default:
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(status)
	bts, err0 := json.Marshal(err)
	if err0 != nil {
		log.Err(err0)
	}
	_, err0 = response.Write(bts)
	if err0 != nil {
		log.Err(err0)
	}
}
