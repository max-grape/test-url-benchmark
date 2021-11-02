package sites

type Benchmark interface {
	Run(url string, limit uint64) (int64, error)
}
