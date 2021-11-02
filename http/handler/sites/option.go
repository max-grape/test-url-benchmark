package sites

type Option func(w *Handler)

func WithParser(parser Parser) Option {
	return func(h *Handler) {
		h.parser = parser
	}
}

func WithBenchmark(benchmark Benchmark) Option {
	return func(h *Handler) {
		h.benchmark = benchmark
	}
}

func WithBenchmarkConcurrency(concurrency uint64) Option {
	return func(h *Handler) {
		h.benchmarkConcurrency = concurrency
	}
}

func WithBenchmarkRequestsLimit(limit uint64) Option {
	return func(h *Handler) {
		h.benchmarkRequestsLimit = limit
	}
}
