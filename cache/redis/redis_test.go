package redis_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	redis2 "github.com/max-grape/test-revo/cache/redis"
	"github.com/stretchr/testify/assert"
)

func mockRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	redisServerMock, err := miniredis.Run()
	assert.NoError(t, err)

	return redisServerMock, redis.NewClient(&redis.Options{
		Addr: redisServerMock.Addr(),
	})
}

func TestNewAbstraction(t *testing.T) {
	_, redisClient := mockRedis(t)

	expected := &redis2.Cache{Client: redisClient}
	actual := redis2.NewAbstraction(redisClient)

	assert.Equal(t, expected, actual)
}

func TestCache_SetNX(t *testing.T) {
	_, redisClient := mockRedis(t)

	redisCache := redis2.NewAbstraction(redisClient)

	result, err := redisCache.SetNX("foo", "bar", time.Second)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestCache_SetNX_Error(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	redisServer.SetError("some error")

	redisCache := redis2.NewAbstraction(redisClient)

	result, err := redisCache.SetNX("foo", "bar", time.Second)
	assert.EqualError(t, err, "some error")
	assert.False(t, result)
}

func TestCache_Set(t *testing.T) {
	_, redisClient := mockRedis(t)

	redisCache := redis2.NewAbstraction(redisClient)

	err := redisCache.Set("foo", "bar", time.Second)
	assert.NoError(t, err)

	actual, err := redisCache.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", actual)
}

func TestCache_Set_Error(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	redisServer.SetError("some error")

	redisCache := redis2.NewAbstraction(redisClient)

	err := redisCache.Set("foo", "bar", time.Second)
	assert.EqualError(t, err, "some error")
}

func TestCache_Get(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	err := redisServer.Set("foo", "bar")
	assert.NoError(t, err)

	redisCache := redis2.NewAbstraction(redisClient)

	actual, err := redisCache.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", actual)
}

func TestCache_Get_Error(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	redisServer.SetError("some error")

	redisCache := redis2.NewAbstraction(redisClient)

	actual, err := redisCache.Get("foo")
	assert.EqualError(t, err, "some error")
	assert.Equal(t, "", actual)
}

func TestCache_Del(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	err := redisServer.Set("foo", "bar")
	assert.NoError(t, err)

	redisCache := redis2.NewAbstraction(redisClient)

	actual, err := redisCache.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", actual)

	err = redisCache.Del("foo")
	assert.NoError(t, err)

	actual, err = redisCache.Get("foo")
	assert.Equal(t, redis.Nil, err)
	assert.Equal(t, "", actual)
}

func TestCache_Del_Error(t *testing.T) {
	redisServer, redisClient := mockRedis(t)

	redisServer.SetError("some error")

	redisCache := redis2.NewAbstraction(redisClient)

	err := redisCache.Del("foo")
	assert.EqualError(t, err, "some error")
}
