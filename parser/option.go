package parser

type Option func(p *parser)

func WithEndpoint(endpoint string) Option {
	return func(p *parser) {
		p.endpoint = endpoint
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(p *parser) {
		p.httpClient = client
	}
}
