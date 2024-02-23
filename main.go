package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/mblancoa/authentication/adapters/cache"
	"github.com/mblancoa/authentication/adapters/mongodb"
	"github.com/mblancoa/authentication/api"
	"github.com/mblancoa/authentication/core"
)

func main() {
	mongodb.SetupMongodbConfiguration()
	cache.SetupRedisCacheConfiguration()
	core.SetupCoreConfig()
	api.SetupApiConfiguration()
	beego.Run()
}
