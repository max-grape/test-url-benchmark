package parser_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/max-grape/test-revo/omap"
	"github.com/max-grape/test-revo/parser"
	"github.com/stretchr/testify/assert"
)

type httpClientMock struct {
	response *http.Response
	err      error
}

func (m *httpClientMock) Get(url string) (response *http.Response, err error) {
	return m.response, m.err
}

func TestParser_Run_Error(t *testing.T) {
	p := parser.New(
		parser.WithEndpoint("endpoint"),
		parser.WithHTTPClient(&httpClientMock{
			err: errors.New("some error"),
		}),
	)

	actual, err := p.Run("foo")

	assert.EqualError(t, err, "failed to do request: some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestParser_Run_InvalidCode(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`empty`))

	p := parser.New(
		parser.WithEndpoint("endpoint"),
		parser.WithHTTPClient(&httpClientMock{
			response: &http.Response{StatusCode: http.StatusInternalServerError, Status: "500", Body: body},
		}),
	)

	actual, err := p.Run("foo")

	assert.EqualError(t, err, "invalid response status code: 500")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestParser_Run(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`<div class="serp-item" data-cid=1>
	<a target="_blank" class="Link" href="http://external-stub-1:8080/some/path"></a>
</div>
<div class="serp-item" data-cid=1>
	<a target="_blank" class="Link" href="http://external-stub-2:8080/some/path"></a>
</div>
`))

	p := parser.New(
		parser.WithEndpoint("endpoint"),
		parser.WithHTTPClient(&httpClientMock{
			response: &http.Response{StatusCode: http.StatusOK, Body: body},
		}),
	)

	actual, err := p.Run("foo")

	assert.NoError(t, err)
	assert.Equal(t, omap.OrderedMap{
		omap.KeyVal{Key: "external-stub-1:8080", Val: "http://external-stub-1:8080/some/path"},
		omap.KeyVal{Key: "external-stub-2:8080", Val: "http://external-stub-2:8080/some/path"},
	}, actual)
}
