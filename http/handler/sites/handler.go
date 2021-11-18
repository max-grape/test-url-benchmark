package sites

import (
	"net/http"
	"sync"
	"time"

	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
	"github.com/max-grape/test-url-benchmark/omap"
)

type ParserWrapper interface {
	Run(searchParam string, fn func() (omap.OrderedMap, error))
}

const (
	defaultBenchmarkConcurrency   = 5
	defaultBenchmarkRequestsLimit = 100

	msgSearchParamIsNotSpecified = "search parameter is not specified"
)

type Handler struct {
	parser                 Parser
	benchmark              Benchmark
	benchmarkConcurrency   uint64
	benchmarkRequestsLimit uint64
}

func New(options ...Option) *Handler {
	h := &Handler{
		benchmarkConcurrency:   defaultBenchmarkConcurrency,
		benchmarkRequestsLimit: defaultBenchmarkRequestsLimit,
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		startedAt = time.Now()
		entry     = ctxlogrus.Extract(r.Context())
	)

	searchParam := r.URL.Query().Get("search")
	if searchParam == "" {
		entry.Error(msgSearchParamIsNotSpecified)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	parseResults, err := h.parser.Run(searchParam)
	if err != nil {
		entry.Errorf("failed to parse `%s` search request: %+v", searchParam, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make(chan item)

	go func() {
		for i, pair := range parseResults {
			url, ok := pair.Val.(string)
			if !ok {
				entry.Errorf("failed to cast `%+v` to string", pair.Val)
				continue
			}

			items <- item{
				i:    i,
				host: pair.Key,
				url:  url,
			}
		}

		close(items)
	}()

	var (
		wg   = sync.WaitGroup{}
		data = make(omap.OrderedMap, len(parseResults))
	)

	for i := uint64(0); i < h.benchmarkConcurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for item := range items {
				result, err := h.benchmark.Run(item.url, h.benchmarkRequestsLimit)
				if err == nil {
					entry.Debugf("benchmarking `%s` url is done; hits: %d", item.url, result)
				} else {
					entry.Debugf("benchmarking `%s` url is done; hits: %d; error: %+v", item.url, result, err)
				}

				data[item.i] = omap.KeyVal{
					Key: item.host,
					Val: result,
				}
			}
		}()
	}

	wg.Wait()

	body, err := data.MarshalIndent("", "    ")
	if err != nil {
		entry.Errorf("failed to marshal response: %+v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

	if _, err := rw.Write(body); err != nil {
		entry.Errorf("failed to write data: %+v", err)
	}

	entry.Debugf("time to process request = %s", time.Since(startedAt))
}

type item struct {
	i    int
	host string
	url  string
}
