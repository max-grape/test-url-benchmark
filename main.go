package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-redis/redis"
	"github.com/improbable-eng/go-httpwares"
	http_logrus "github.com/improbable-eng/go-httpwares/logging/logrus"
	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
	"github.com/max-grape/test-revo/benchmark"
	cacheRedis "github.com/max-grape/test-revo/cache/redis"
	"github.com/max-grape/test-revo/http/handler/health"
	"github.com/max-grape/test-revo/http/handler/sites"
	"github.com/max-grape/test-revo/parser"
	"github.com/max-grape/test-revo/shutdown"
	benchmarkCacheWrapper "github.com/max-grape/test-revo/wrapper/benchmark"
	parserCacheWrapper "github.com/max-grape/test-revo/wrapper/parser"
	"github.com/sirupsen/logrus"
)

func main() {
	// Config

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("failed to get config: %+v", err)
	}

	// Logger

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse logger level: %+v", err)
	}

	logger := logrus.New()

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	entry := logrus.NewEntry(logger)

	entry.Info("starting service...")

	entry.Debugf("config = %+v", cfg)

	// CPU

	runtime.GOMAXPROCS(cfg.CPU)

	// Recovery

	defer func() {
		if p := recover(); p != nil {
			entry.Errorf("panic recovered: %+v; stack trace: %s", p, string(debug.Stack()))
		}
	}()

	// Shutdowner

	shutdowner := shutdown.New(
		shutdown.WithTimeout(cfg.ShutdownTimeout),
	)

	// Health HTTP handler

	healthHandler := health.New()

	// Redis client

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	healthHandler.RegisterHealthCheckFunc(func() bool {
		if _, err := redisClient.Ping().Result(); err != nil {
			entry.Errorf("failed to ping Redis: %+v", err)
			return false
		}

		return true
	})

	shutdowner.RegisterShutdownFunc(func() {
		if err := redisClient.Close(); err != nil {
			entry.Errorf("failed to close Redis client: %+v", err)
		}
	})

	// Redis abstraction

	cache := cacheRedis.NewAbstraction(redisClient)

	// HTTP client

	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: net.IPv4zero},
		KeepAlive: time.Second * 30,
	}

	httpClient := &http.Client{
		Timeout: cfg.HTTPClientTimeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                dialer.Dial,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: 10000,
			MaxConnsPerHost:     0,
		},
	}

	// Parser

	p := parser.New(
		parser.WithEndpoint(cfg.ParserEndpoint),
		parser.WithHTTPClient(httpClient),
	)

	// Parser cache layer

	p = parserCacheWrapper.New(
		parserCacheWrapper.WithParser(p),
		parserCacheWrapper.WithCache(cache),
		parserCacheWrapper.WithInProgressCacheExpiration(cfg.ParserInProgressCacheExpiration),
		parserCacheWrapper.WithCacheExpiration(cfg.ParserCacheExpiration),
	)

	// Benchmark

	b := benchmark.New(
		benchmark.WithHTTPClient(httpClient),
	)

	// Benchmark cache layer

	b = benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(b),
		benchmarkCacheWrapper.WithCache(cache),
		benchmarkCacheWrapper.WithInProgressCacheExpiration(cfg.BenchmarkInProgressCacheExpiration),
		benchmarkCacheWrapper.WithCacheExpiration(cfg.BenchmarkCacheExpiration),
	)

	// Sites handler

	sitesHandler := sites.New(
		sites.WithParser(p),
		sites.WithBenchmark(b),
		sites.WithBenchmarkConcurrency(cfg.BenchmarkConcurrency),
		sites.WithBenchmarkRequestsLimit(cfg.BenchmarkRequestsLimit),
	)

	// Logger middleware

	loggerMiddleware := http_logrus.Middleware(entry,
		http_logrus.WithDecider(func(w httpwares.WrappedResponseWriter, r *http.Request) bool {
			return w.StatusCode() != http.StatusOK
		}),
		http_logrus.WithRequestFieldExtractor(func(r *http.Request) map[string]interface{} {
			fields := make(map[string]interface{})

			switch r.URL.Path {
			case "/sites":
				fields["search"] = r.URL.Query().Get("search")
			}

			return fields
		}),
	)

	// Recovery middleware

	recoveryMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					entry := ctxlogrus.Extract(r.Context())
					entry.Errorf("panic recovered: %+v; stack trace: %s", err, string(debug.Stack()))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}

	// External HTTP server

	externalMux := http.NewServeMux()

	externalMux.Handle("/sites", loggerMiddleware(recoveryMiddleware(sitesHandler)))

	externalServer := &http.Server{
		Addr:         cfg.HTTPExternalServerAddress,
		Handler:      externalMux,
		ReadTimeout:  cfg.HTTPExternalServerReadTimeout,
		WriteTimeout: cfg.HTTPExternalServerWriteTimeout,
	}

	go func() {
		if err := externalServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			entry.Errorf("failed to listen and serve external HTTP server: %+v", err)
		}
	}()

	shutdowner.RegisterShutdownFunc(func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPExternalServerShutdownTimeout)
		defer cancel()

		if err := externalServer.Shutdown(ctx); err != nil {
			entry.Errorf("failed to shutdown external HTTP server: %+v", err)
		}
	})

	// HTTP internal server

	internalMux := http.NewServeMux()

	internalMux.Handle("/health", loggerMiddleware(recoveryMiddleware(healthHandler)))

	internalServer := &http.Server{
		Addr:         cfg.HTTPInternalServerAddress,
		Handler:      internalMux,
		ReadTimeout:  cfg.HTTPInternalServerReadTimeout,
		WriteTimeout: cfg.HTTPInternalServerWriteTimeout,
	}

	go func() {
		if err := internalServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			entry.Errorf("failed to listen and serve internal HTTP server: %+v", err)
		}
	}()

	shutdowner.RegisterShutdownFunc(func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPInternalServerShutdownTimeout)
		defer cancel()

		if err := internalServer.Shutdown(ctx); err != nil {
			entry.Errorf("failed to shutdown internal HTTP server: %+v", err)
		}
	})

	// Waiting for shutdown

	entry.Info("service is started")

	if err := shutdowner.WaitForShutdown(); err != nil {
		entry.Errorf("failed to shutdown: %+v", err)
	}

	entry.Info("service is stopped")
}
