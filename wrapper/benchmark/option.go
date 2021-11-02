package benchmark

import (
	"time"
)

type Option func(w *Wrapper)

func WithBenchmark(benchmark Benchmark) Option {
	return func(w *Wrapper) {
		w.benchmark = benchmark
	}
}

func WithCache(cache Cache) Option {
	return func(w *Wrapper) {
		w.cache = cache
	}
}

func WithInProgressCacheExpiration(time time.Duration) Option {
	return func(w *Wrapper) {
		w.inProgressCacheExpiration = time
	}
}

func WithCacheExpiration(time time.Duration) Option {
	return func(w *Wrapper) {
		w.cacheExpiration = time
	}
}
