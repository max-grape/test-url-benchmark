package benchmark_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/max-grape/test-url-benchmark/benchmark"
	"github.com/stretchr/testify/assert"
)

type httpClientMock struct {
	response func() *http.Response
	err      error
}

func (m *httpClientMock) Get(url string) (*http.Response, error) {
	var response *http.Response

	if m.response != nil {
		response = m.response()
	}

	return response, m.err
}

func TestBenchmark_Run_Error(t *testing.T) {
	var count uint64

	httpClient := &httpClientMock{
		response: func() *http.Response {
			atomic.AddUint64(&count, 1)
			return nil
		},
		err: errors.New("some error"),
	}

	b := benchmark.New(
		benchmark.WithHTTPClient(httpClient),
	)

	actual, err := b.Run("foo", 10)

	assert.EqualError(t, err, "some error")
	assert.Equal(t, int64(0), actual)
	assert.Equal(t, uint64(1), count)
}

func TestBenchmark_Run_ResponseCodeIsBelow200(t *testing.T) {
	var (
		count uint64
		body  = ioutil.NopCloser(bytes.NewBufferString(`empty`))
	)

	httpClient := &httpClientMock{
		response: func() *http.Response {
			atomic.AddUint64(&count, 1)
			return &http.Response{StatusCode: 199, Status: "199", Body: body}
		},
	}

	b := benchmark.New(
		benchmark.WithHTTPClient(httpClient),
	)

	actual, err := b.Run("foo", 10)

	assert.EqualError(t, err, "199")
	assert.Equal(t, int64(0), actual)
	assert.Equal(t, uint64(1), count)
}

func TestBenchmark_Run_ResponseCodeIsAbove399(t *testing.T) {
	var (
		count uint64
		body  = ioutil.NopCloser(bytes.NewBufferString(`empty`))
	)

	httpClient := &httpClientMock{
		response: func() *http.Response {
			atomic.AddUint64(&count, 1)
			return &http.Response{StatusCode: http.StatusBadRequest, Status: "400", Body: body}
		},
	}

	b := benchmark.New(
		benchmark.WithHTTPClient(httpClient),
	)

	actual, err := b.Run("foo", 10)

	assert.EqualError(t, err, "400")
	assert.Equal(t, int64(0), actual)
	assert.Equal(t, uint64(1), count)
}

func TestBenchmark_Run_Success(t *testing.T) {
	var (
		count uint64
		body  = ioutil.NopCloser(bytes.NewBufferString(`empty`))
	)

	httpClient := &httpClientMock{
		response: func() *http.Response {
			atomic.AddUint64(&count, 1)
			return &http.Response{StatusCode: http.StatusOK, Status: "200", Body: body}
		},
	}

	b := benchmark.New(
		benchmark.WithHTTPClient(httpClient),
	)

	actual, err := b.Run("foo", 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), actual)
	assert.Equal(t, uint64(10), count)
}
