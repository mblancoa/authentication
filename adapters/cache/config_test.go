package cache

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/mblancoa/authentication/core"
	"github.com/mblancoa/authentication/tools"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func init() {
	err := os.Chdir("./../..")
	tools.ManageTestError(err)
	err = os.Setenv(core.RunMode, "test")
	tools.ManageTestError(err)
}

func TestLoadConfiguration(t *testing.T) {
	var config configuration
	core.LoadYamlConfiguration(core.GetConfigFile(), &config)

	opt := options{Addr: "localhost:16379", Password: "", DB: 0}
	timeout := time.Duration(5) * time.Minute
	assert.NotEmpty(t, config)
	assert.NotEmpty(t, config.Redis)

	checkEmail := config.Redis.CheckEmailCache
	assert.NotEmpty(t, checkEmail)
	assert.Equal(t, opt, checkEmail.Options)
	assert.Equal(t, "test-check-email:*", checkEmail.KeyPattern)
	assert.Equal(t, timeout, checkEmail.Timeout)

	codeConfirmation := config.Redis.CodeConfirmationCache
	assert.NotEmpty(t, checkEmail)
	assert.Equal(t, opt, codeConfirmation.Options)
	assert.Equal(t, "test-code-confirmation:*", codeConfirmation.KeyPattern)
	assert.Equal(t, timeout, codeConfirmation.Timeout)
}

func TestSetupRedisCacheConfiguration(t *testing.T) {
	mini := miniredis.NewMiniRedis()
	err := mini.StartAddr("localhost:16379")
	tools.ManageTestError(err)
	defer mini.Close()

	SetupRedisCacheConfiguration()

	assert.NotEmpty(t, core.CacheContext)
	assert.NotEmpty(t, core.CacheContext.CheckEmailCache)
	assert.NotEmpty(t, core.CacheContext.CodeConfirmationCache)
}
