package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/mblancoa/authentication/api"
	"github.com/mblancoa/authentication/config"
	_ "github.com/mblancoa/authentication/config"
)

func init() {
	config.SetupMongodbConfiguration()
	config.SetupCoreConfiguration()
	config.SetupApiConfiguration()
}

func main() {
	beego.Run()
}
