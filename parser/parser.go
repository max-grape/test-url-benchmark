package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/max-grape/test-url-benchmark/omap"
)

const (
	defaultEndpoint = "https://yandex.ru/search/touch"
)

type Parser interface {
	Run(searchRequest string) (omap.OrderedMap, error)
}

type parser struct {
	endpoint   string
	httpClient HTTPClient
}

func New(options ...Option) Parser {
	p := &parser{
		endpoint: defaultEndpoint,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *parser) Run(searchRequest string) (omap.OrderedMap, error) {
	searchURL := fmt.Sprintf(baseYandexURL, p.endpoint, url.QueryEscape(searchRequest))

	response, err := p.httpClient.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status code: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	parseResponse := parseYandexResponse(body)

	if parseResponse.Error != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", parseResponse.Error)
	}

	result := make(omap.OrderedMap, len(parseResponse.Items))

	for i, item := range parseResponse.Items {
		result[i] = omap.KeyVal{
			Key: item.Host,
			Val: item.Url,
		}
	}

	return result, nil
}
