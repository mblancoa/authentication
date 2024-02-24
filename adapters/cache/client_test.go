package cache

import (
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/pioz/faker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var testKeyPattern string = "test:*"
var testKey string = "testKey"
var testCacheKey string = "test:testKey"

type redisCacheSuite struct {
	suite.Suite
	mini  *miniredis.Miniredis
	cache *redisClient
}

func TestRedisCacheSuite(t *testing.T) {
	suite.Run(t, new(redisCacheSuite))
}

func (suite *redisCacheSuite) SetupSuite() {
	suite.mini, _ = miniredis.Run()
	opt := &redis.Options{Addr: suite.mini.Addr(), Password: "", DB: 0}
	suite.cache = &redisClient{Client: redis.NewClient(opt), timeout: time.Minute, keyPattern: testKeyPattern}
}

func (suite *redisCacheSuite) SetupTest() {
}

func (suite *redisCacheSuite) TearDownTest() {
	suite.mini.Del(testCacheKey)
}

func (suite *redisCacheSuite) TearDownSuite() {
	suite.mini.Close()
}

func (suite *redisCacheSuite) TestSetString() {

	s := faker.StringWithSize(24)

	err := suite.cache.Set(testKey, s)

	suite.Assertions.NoError(err)
	inRedis, _ := suite.mini.Get(testCacheKey)
	suite.Assertions.Equal(s, inRedis)
}

func (suite *redisCacheSuite) TestSetNil() {

	err := suite.cache.Set(testKey, nil)

	suite.Assertions.NoError(err)
	inRedis, _ := suite.mini.Get(testCacheKey)
	suite.Assertions.Equal("", inRedis)
}

func (suite *redisCacheSuite) TestSetObject() {

	obj := struct {
		Prop1, Prop2 string
		Prop3        struct{ Prop1, Prop2 string }
	}{
		Prop1: "test1", Prop2: "test2",
		Prop3: struct{ Prop1, Prop2 string }{Prop1: "sub1", Prop2: "sub2"},
	}

	bts, _ := json.Marshal(obj)
	expected := string(bts)

	err := suite.cache.Set(testKey, obj)

	suite.Assertions.NoError(err)
	inRedis, _ := suite.mini.Get(testCacheKey)
	suite.Assertions.Equal(expected, inRedis)
}

func (suite *redisCacheSuite) TestGet() {
	type testStruct struct {
		Prop1, Prop2 string
		Prop3        struct{ Prop1, Prop2 string }
	}
	obj := testStruct{
		Prop1: "test1", Prop2: "test2",
		Prop3: struct{ Prop1, Prop2 string }{Prop1: "sub1", Prop2: "sub2"},
	}

	bts, _ := json.Marshal(obj)
	inRedis := string(bts)
	_ = suite.mini.Set(testCacheKey, inRedis)

	out := testStruct{}
	err := suite.cache.Get(testKey, &out)

	suite.Assertions.NoError(err)
	suite.Assertions.Equal(obj, out)
	suite.Assertions.True(suite.mini.Exists(testCacheKey))
}

func (suite *redisCacheSuite) TestGetAndDelete() {
	type testStruct struct {
		Prop1, Prop2 string
		Prop3        struct{ Prop1, Prop2 string }
	}
	obj := testStruct{
		Prop1: "test1", Prop2: "test2",
		Prop3: struct{ Prop1, Prop2 string }{Prop1: "sub1", Prop2: "sub2"},
	}

	bts, _ := json.Marshal(obj)
	inRedis := string(bts)
	_ = suite.mini.Set(testCacheKey, inRedis)

	out := testStruct{}
	err := suite.cache.GetAndDelete(testKey, &out)

	suite.Assertions.NoError(err)
	suite.Assertions.Equal(obj, out)
	suite.Assertions.False(suite.mini.Exists(testCacheKey))
}

func (suite *redisCacheSuite) TestGetString() {
	str := faker.StringWithSize(24)
	_ = suite.mini.Set(testCacheKey, str)

	out, err := suite.cache.GetString(testKey)

	suite.Assertions.NoError(err)
	suite.Assertions.Equal(str, out)
	suite.Assertions.True(suite.mini.Exists(testCacheKey))
}

func (suite *redisCacheSuite) TestGetStringAndDelete() {
	str := faker.StringWithSize(24)
	_ = suite.mini.Set(testCacheKey, str)

	out, err := suite.cache.GetStringAndDelete(testKey)

	suite.Assertions.NoError(err)
	suite.Assertions.Equal(str, out)
	suite.Assertions.False(suite.mini.Exists(testCacheKey))
}
