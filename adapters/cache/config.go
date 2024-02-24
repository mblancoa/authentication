package cache

import (
	"context"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/core/ports"
	"github.com/mblancoa/authentication/errors"
	"github.com/mblancoa/authentication/tools"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type options struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type cacheConfig struct {
	Options    options       `yaml:"options"`
	Timeout    time.Duration `yaml:"timeout"`
	KeyPattern string        `yaml:"key_pattern"`
}
type configuration struct {
	Redis struct {
		CheckEmailCache       cacheConfig `yaml:"check_email"`
		CodeConfirmationCache cacheConfig `yaml:"code_confirmation"`
	} `yaml:"redis"`
}

func SetupRedisCacheConfiguration() {
	log.Info().Msg("Initializing redis configuration")
	var config configuration
	core.LoadYamlConfiguration(core.GetConfigFile(), &config)

	setupCacheContext(config)
}

func setupCacheContext(config configuration) {
	core.CacheContext.CheckEmailCache = newCache(config.Redis.CheckEmailCache)
	core.CacheContext.CodeConfirmationCache = newCache(config.Redis.CodeConfirmationCache)
}

func newCache(c cacheConfig) ports.Cache {
	opt := &redis.Options{}
	tools.Mapper(&c.Options, opt)
	client := redis.NewClient(opt)
	_, err := client.Ping(context.Background()).Result()
	errors.ManageErrorPanic(err)
	return NewCache(c.KeyPattern, c.Timeout, client)
}
