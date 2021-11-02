package parser

import (
	"time"
)

type Cache interface {
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(keys ...string) error
}
