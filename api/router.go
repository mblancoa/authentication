package api

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/mblancoa/authentication/api/controllers"
)

func init() {
	beego.Router("/login", &controllers.AuthenticationController{}, "post:Login")
}
