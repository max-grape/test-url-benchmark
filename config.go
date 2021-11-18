package main

import (
	"runtime"
	"time"

	"github.com/max-grape/config/env"
)

const (
	defaultLogLevel = "info"

	defaultRedisHost     = "redis:6379"
	defaultRedisPassword = ""
	defaultRedisDatabase = 0

	defaultHTTPClientTimeout             = time.Second * 3
	defaultHTTPClientKeepAlive           = time.Second * 30
	defaultHTTPClientMaxIdleConnsPerHost = 10000
	defaultHTTPClientMaxConnsPerHost     = 0

	defaultParserEndpoint                  = ""
	defaultParserInProgressCacheExpiration = time.Second * 4
	defaultParserCacheExpiration           = time.Hour * 24

	defaultBenchmarkConcurrency               = 5
	defaultBenchmarkRequestsLimit             = 100
	defaultBenchmarkInProgressCacheExpiration = time.Second * 4
	defaultBenchmarkCacheExpiration           = time.Hour * 24

	defaultHTTPExternalServerAddress         = ":8080"
	defaultHTTPExternalServerReadTimeout     = time.Second * 2
	defaultHTTPExternalServerWriteTimeout    = time.Second * 2
	defaultHTTPExternalServerShutdownTimeout = time.Second * 10

	defaultHTTPInternalServerAddress         = ":8090"
	defaultHTTPInternalServerReadTimeout     = time.Second * 2
	defaultHTTPInternalServerWriteTimeout    = time.Second * 2
	defaultHTTPInternalServerShutdownTimeout = time.Second * 10

	defaultShutdownTimeout = time.Second * 10
)

type config struct {
	LogLevel string `env:"LOG_LEVEL"`

	CPU int `env:"CPU"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDatabase int    `env:"REDIS_DATABASE"`

	HTTPClientTimeout             time.Duration `env:"HTTP_CLIENT_TIMEOUT"`
	HTTPClientKeepAlive           time.Duration `env:"HTTP_CLIENT_KEEP_ALIVE"`
	HTTPClientMaxIdleConnsPerHost int           `env:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"`
	HTTPClientMaxConnsPerHost     int           `env:"HTTP_CLIENT_MAX_CONNS_PER_HOST"`

	ParserEndpoint                  string        `env:"PARSER_ENDPOINT"`
	ParserInProgressCacheExpiration time.Duration `env:"PARSER_IN_PROGRESS_CACHE_EXPIRATION"`
	ParserCacheExpiration           time.Duration `env:"PARSER_CACHE_EXPIRATION"`

	BenchmarkConcurrency               uint64        `env:"BENCHMARK_CONCURRENCY"`
	BenchmarkRequestsLimit             uint64        `env:"BENCHMARK_REQUESTS_LIMIT"`
	BenchmarkInProgressCacheExpiration time.Duration `env:"BENCHMARK_IN_PROGRESS_CACHE_EXPIRATION"`
	BenchmarkCacheExpiration           time.Duration `env:"BENCHMARK_CACHE_EXPIRATION"`

	HTTPExternalServerAddress         string        `env:"HTTP_EXTERNAL_SERVER_ADDRESS"`
	HTTPExternalServerReadTimeout     time.Duration `env:"HTTP_EXTERNAL_SERVER_READ_TIMEOUT"`
	HTTPExternalServerWriteTimeout    time.Duration `env:"HTTP_EXTERNAL_SERVER_WRITE_TIMEOUT"`
	HTTPExternalServerShutdownTimeout time.Duration `env:"HTTP_EXTERNAL_SERVER_SHUTDOWN_TIMEOUT"`

	HTTPInternalServerAddress         string        `env:"HTTP_INTERNAL_SERVER_ADDRESS"`
	HTTPInternalServerReadTimeout     time.Duration `env:"HTTP_INTERNAL_SERVER_READ_TIMEOUT"`
	HTTPInternalServerWriteTimeout    time.Duration `env:"HTTP_INTERNAL_SERVER_WRITE_TIMEOUT"`
	HTTPInternalServerShutdownTimeout time.Duration `env:"HTTP_INTERNAL_SERVER_SHUTDOWN_TIMEOUT"`

	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

func getConfig() (*config, error) {
	cfg := &config{
		LogLevel: defaultLogLevel,

		CPU: runtime.NumCPU(),

		RedisHost:     defaultRedisHost,
		RedisPassword: defaultRedisPassword,
		RedisDatabase: defaultRedisDatabase,

		HTTPClientTimeout:             defaultHTTPClientTimeout,
		HTTPClientKeepAlive:           defaultHTTPClientKeepAlive,
		HTTPClientMaxIdleConnsPerHost: defaultHTTPClientMaxIdleConnsPerHost,
		HTTPClientMaxConnsPerHost:     defaultHTTPClientMaxConnsPerHost,

		ParserEndpoint:                  defaultParserEndpoint,
		ParserInProgressCacheExpiration: defaultParserInProgressCacheExpiration,
		ParserCacheExpiration:           defaultParserCacheExpiration,

		BenchmarkConcurrency:               defaultBenchmarkConcurrency,
		BenchmarkRequestsLimit:             defaultBenchmarkRequestsLimit,
		BenchmarkInProgressCacheExpiration: defaultBenchmarkInProgressCacheExpiration,
		BenchmarkCacheExpiration:           defaultBenchmarkCacheExpiration,

		HTTPExternalServerAddress:         defaultHTTPExternalServerAddress,
		HTTPExternalServerReadTimeout:     defaultHTTPExternalServerReadTimeout,
		HTTPExternalServerWriteTimeout:    defaultHTTPExternalServerWriteTimeout,
		HTTPExternalServerShutdownTimeout: defaultHTTPExternalServerShutdownTimeout,

		HTTPInternalServerAddress:         defaultHTTPInternalServerAddress,
		HTTPInternalServerReadTimeout:     defaultHTTPInternalServerReadTimeout,
		HTTPInternalServerWriteTimeout:    defaultHTTPInternalServerWriteTimeout,
		HTTPInternalServerShutdownTimeout: defaultHTTPInternalServerShutdownTimeout,

		ShutdownTimeout: defaultShutdownTimeout,
	}

	return cfg, env.Parse(cfg)
}
