package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
	"net/http"
)

type AuthenticationController struct {
	beego.Controller
	AuthenticationController core.AuthenticationService
}

func (c *AuthenticationController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"

}

func (c *AuthenticationController) Login() {
	var rq LoginRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &rq)

	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}

	credentials := core.Credentials{}
	tools.Mapper(&rq, &credentials)

	authenticationService := WebApplicationContext.AuthenticationService

	jwt, err := authenticationService.Login(credentials)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	err = c.Ctx.JSONResp(LoginResponse{Token: jwt})
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
	}
}
