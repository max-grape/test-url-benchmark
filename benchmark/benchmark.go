package benchmark

import (
	"errors"
	"sync"
)

type Benchmark interface {
	Run(url string, limit uint64) (int64, error)
}

type benchmark struct {
	httpClient HTTPClient
}

func New(options ...Option) Benchmark {
	b := &benchmark{}

	for _, option := range options {
		option(b)
	}

	return b
}

func (b *benchmark) Run(url string, limit uint64) (int64, error) {
	var (
		runners = uint64(0)
		wg      = sync.WaitGroup{}
		jobs    = make(chan struct{})
		results = make(chan error)
		stop    = make(chan struct{})
	)

	go func() {
		defer func() {
			close(jobs)
			wg.Wait()
			close(results)
		}()

		count := uint64(0)

		for {
			if count >= limit {
				return
			}

			if runners < limit {
				select {
				case jobs <- struct{}{}:
					count++
					continue
				case <-stop:
					return
				default:
					runners++
					wg.Add(1)
					go b.runner(url, &wg, jobs, results)
				}
			}

			select {
			case jobs <- struct{}{}:
				count++
			case <-stop:
				return
			}
		}
	}()

	var (
		result int64
		err    error
	)

	for err = range results {
		if err != nil {
			close(stop)
			break
		} else {
			result++
		}
	}

	return result, err
}

func (b *benchmark) runner(url string, wg *sync.WaitGroup, jobs <-chan struct{}, results chan<- error) {
	defer wg.Done()

	for range jobs {
		results <- b.request(url)
	}
}

func (b *benchmark) request(url string) error {
	response, err := b.httpClient.Get(url)
	if err != nil {
		return err
	}

	_ = response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return errors.New(response.Status)
	}

	return nil
}
