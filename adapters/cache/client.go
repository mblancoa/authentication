package cache

import (
	"context"
	"encoding/json"
	"github.com/mblancoa/authentication/core/ports"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"time"
)

type redisClient struct {
	*redis.Client
	keyPattern string
	timeout    time.Duration
}

func NewCache(keyPattern string, timeout time.Duration, client *redis.Client) ports.Cache {
	return &redisClient{
		Client:     client,
		keyPattern: keyPattern,
		timeout:    timeout,
	}
}

func (r *redisClient) Set(key string, v interface{}) error {
	k := r.generateKey(key)
	if v == nil || reflect.TypeOf(v).Name() == "string" {
		return r.Client.Set(context.Background(), k, v, r.timeout).Err()
	}
	bts, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.Client.Set(context.Background(), k, string(bts), r.timeout).Err()
}

func (r *redisClient) Get(key string, v interface{}) error {
	k := r.generateKey(key)
	bts, err := r.Client.Get(context.Background(), k).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, v)
}

func (r *redisClient) GetAndDelete(key string, v interface{}) error {
	k := r.generateKey(key)
	bts, err := r.Client.GetDel(context.Background(), k).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, v)
}

func (r *redisClient) GetString(key string) (string, error) {
	k := r.generateKey(key)
	return r.Client.Get(context.Background(), k).Result()
}

func (r *redisClient) GetStringAndDelete(key string) (string, error) {
	k := r.generateKey(key)
	return r.Client.GetDel(context.Background(), k).Result()
}

func (r *redisClient) generateKey(key string) string {
	return strings.Replace(r.keyPattern, "*", key, 1)
}
