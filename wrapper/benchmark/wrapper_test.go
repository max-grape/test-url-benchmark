package benchmark_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	benchmarkCacheWrapper "github.com/max-grape/test-revo/wrapper/benchmark"
	"github.com/stretchr/testify/assert"
)

type benchmarkMock struct {
	data int64
	err  error
}

func (b *benchmarkMock) Run(url string, limit uint64) (int64, error) {
	return b.data, b.err
}

type cacheMock struct {
	setNxResponse bool
	setNxErr      error
	setErr        error
	getResponse   string
	getErr        error
	delErr        error
}

func (c *cacheMock) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.setNxResponse, c.setNxErr
}

func (c *cacheMock) Set(key string, value interface{}, expiration time.Duration) error {
	return c.setErr
}

func (c *cacheMock) Get(key string) (string, error) {
	return c.getResponse, c.getErr
}

func (c *cacheMock) Del(keys ...string) error {
	return c.delErr
}

func TestWrapper_Run_SetNXError(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxErr: errors.New("some error"),
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, "failed to set processing `benchmark:foo` key to cache: some error")
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_SetError(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: true,
			setErr:        fmt.Errorf("some error"),
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, "failed to set `benchmark:foo` key to cache: some error")
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_LockedSuccess(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{
			data: 123,
		}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: true,
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), actual)
}

func TestWrapper_Run_GetError(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getErr:        fmt.Errorf("some error"),
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, "failed to get `benchmark:foo` key from cache: some error")
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_GetValueInProgress(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "inprogress",
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, benchmarkCacheWrapper.ErrInProgress.Error())
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_ParseIntError(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "invalid_int64",
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, "failed to parse `invalid_int64` to int64: strconv.ParseInt: parsing \"invalid_int64\": invalid syntax")
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_DelError(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "invalid_int64",
			delErr:        errors.New("some error"),
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.EqualError(t, err, "failed to parse `invalid_int64` to int64: strconv.ParseInt: parsing \"invalid_int64\": invalid syntax; faield to delete cache `benchmark:foo` key: some error")
	assert.Equal(t, int64(0), actual)
}

func TestWrapper_Run_NotLockedSuccess(t *testing.T) {
	pcw := benchmarkCacheWrapper.New(
		benchmarkCacheWrapper.WithBenchmark(&benchmarkMock{}),
		benchmarkCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "123",
		}),
	)

	actual, err := pcw.Run("foo", 100)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), actual)
}
