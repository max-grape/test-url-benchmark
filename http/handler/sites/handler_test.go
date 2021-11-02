package sites

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/max-grape/test-revo/omap"
	"github.com/stretchr/testify/assert"
)

type parserMock struct {
	data func() omap.OrderedMap
	err  error
}

func (p *parserMock) Run(searchParam string) (omap.OrderedMap, error) {
	var data omap.OrderedMap

	if p.data != nil {
		data = p.data()
	}

	return data, p.err
}

type benchmarkMock struct {
	data int64
	err  error
}

func (b *benchmarkMock) Run(url string, limit uint64) (int64, error) {
	return b.data, b.err
}

func TestHandler_ServeHTTP_SearchParameterIsNotSpecified(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sites", nil)

	New().ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_ServeHTTP_SearchParameterIsSpecifiedButEmpty(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sites?search=", nil)

	New().ServeHTTP(rr, r)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_ServeHTTP_ParserRunError(t *testing.T) {
	s := New(
		WithParser(&parserMock{
			err: errors.New("some error"),
		}),
	)

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sites?search=foo", nil)

	s.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_ServeHTTP_BenchmarkRunError(t *testing.T) {
	s := New(
		WithParser(&parserMock{
			data: func() omap.OrderedMap { return omap.OrderedMap{{Key: "host", Val: "url"}} },
		}),
		WithBenchmark(&benchmarkMock{
			err: errors.New("some error"),
		}),
	)

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sites?search=foo", nil)

	s.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "{\n    \"host\": 0\n}", rr.Body.String())
}

func TestHandler_ServeHTTP_Success(t *testing.T) {
	s := New(
		WithParser(&parserMock{
			data: func() omap.OrderedMap { return omap.OrderedMap{{Key: "host", Val: "url"}} },
		}),
		WithBenchmark(&benchmarkMock{
			data: 123,
		}),
	)

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sites?search=foo", nil)

	s.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "{\n    \"host\": 123\n}", rr.Body.String())
}
