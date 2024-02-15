package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/mblancoa/authentication/adapters"
	"github.com/mblancoa/authentication/api"
	"github.com/mblancoa/authentication/core"
)

func main() {
	adapters.SetupMongodbConfiguration()
	core.SetupCoreConfig()
	api.SetupApiConfiguration()
	beego.Run()
}
