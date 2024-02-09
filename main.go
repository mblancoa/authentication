package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/mblancoa/authentication/api"
	"github.com/mblancoa/authentication/config"
)

func main() {
	config.SetupCoreConfiguration()
	config.SetupMongoDBConfiguration()

	beego.Run()
}
