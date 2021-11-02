package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	*redis.Client
}

func NewAbstraction(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(key, value, expiration).Result()
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	_, err := c.Client.Set(key, value, expiration).Result()
	return err
}

func (c *Cache) Get(key string) (string, error) {
	return c.Client.Get(key).Result()
}

func (c *Cache) Del(keys ...string) error {
	_, err := c.Client.Del(keys...).Result()
	return err
}
