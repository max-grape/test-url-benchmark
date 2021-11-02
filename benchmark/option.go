package benchmark

type Option func(b *benchmark)

func WithHTTPClient(client HTTPClient) Option {
	return func(b *benchmark) {
		b.httpClient = client
	}
}
