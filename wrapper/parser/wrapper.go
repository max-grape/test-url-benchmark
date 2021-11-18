package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/max-grape/test-url-benchmark/omap"
)

const (
	cacheKeyPrefix                   = "parser:"
	cacheValueInProgress             = "inprogress"
	defaultInProgressCacheExpiration = time.Second * 4
	defaultCacheExpiration           = time.Hour * 24
)

var (
	ErrInProgress = errors.New("in progress by another request")
)

type Wrapper struct {
	parser                    Parser
	cache                     Cache
	inProgressCacheExpiration time.Duration
	cacheExpiration           time.Duration
}

func New(options ...Option) Parser {
	w := &Wrapper{
		inProgressCacheExpiration: defaultInProgressCacheExpiration,
		cacheExpiration:           defaultCacheExpiration,
	}

	for _, option := range options {
		option(w)
	}

	return w
}

func (w *Wrapper) Run(searchRequest string) (omap.OrderedMap, error) {
	var (
		key    = cacheKeyPrefix + searchRequest
		result omap.OrderedMap
	)

	locked, err := w.cache.SetNX(key, cacheValueInProgress, w.inProgressCacheExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to set processing `%s` key to cache: %w", key, err)
	}

	if locked {
		result, err = w.parser.Run(searchRequest)
		if err != nil {
			return nil, err
		}

		if err := w.cache.Set(key, result, w.cacheExpiration); err != nil {
			return nil, fmt.Errorf("failed to set `%s` key to cache: %w", key, err)
		}
	} else {
		data, err := w.cache.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get `%s` key from cache: %w", key, err)
		}

		if data == cacheValueInProgress {
			return nil, ErrInProgress
		}

		if err := json.Unmarshal([]byte(data), &result); err != nil {
			err = fmt.Errorf("failed to unmarhsal `%s` key: %+v", key, err)

			if delErr := w.cache.Del(key); delErr != nil {
				err = fmt.Errorf("%s; faield to delete cache `%s` key: %s", err, key, delErr)
			}

			return nil, err
		}
	}

	return result, nil
}
