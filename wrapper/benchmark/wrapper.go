package benchmark

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	cacheKeyPrefix                   = "benchmark:"
	cacheValueInProgress             = "inprogress"
	defaultInProgressCacheExpiration = time.Second * 4
	defaultCacheExpiration           = time.Hour * 24
)

var (
	ErrInProgress = errors.New("in progress by another request")
)

type Wrapper struct {
	benchmark                 Benchmark
	cache                     Cache
	inProgressCacheExpiration time.Duration
	cacheExpiration           time.Duration
}

func New(options ...Option) Benchmark {
	w := &Wrapper{
		inProgressCacheExpiration: defaultInProgressCacheExpiration,
		cacheExpiration:           defaultCacheExpiration,
	}

	for _, option := range options {
		option(w)
	}

	return w
}

func (w *Wrapper) Run(url string, limit uint64) (int64, error) {
	var (
		key    = cacheKeyPrefix + url
		result int64
	)

	locked, err := w.cache.SetNX(key, cacheValueInProgress, w.inProgressCacheExpiration)
	if err != nil {
		return 0, fmt.Errorf("failed to set processing `%s` key to cache: %w", key, err)
	}

	if locked {
		result, err = w.benchmark.Run(url, limit)

		if err := w.cache.Set(key, result, w.cacheExpiration); err != nil {
			return 0, fmt.Errorf("failed to set `%s` key to cache: %+v", key, err)
		}
	} else {
		resultStr, err := w.cache.Get(key)
		if err != nil {
			return 0, fmt.Errorf("failed to get `%s` key from cache: %+v", key, err)
		}

		if resultStr == cacheValueInProgress {
			return 0, ErrInProgress
		}

		if result, err = strconv.ParseInt(resultStr, 10, 64); err != nil {
			err = fmt.Errorf("failed to parse `%s` to int64: %+v", resultStr, err)

			if delErr := w.cache.Del(key); delErr != nil {
				err = fmt.Errorf("%s; faield to delete cache `%s` key: %s", err, key, delErr)
			}

			return 0, err
		}
	}

	return result, err
}
