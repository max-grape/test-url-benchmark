package parser_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/max-grape/test-url-benchmark/omap"
	parserCacheWrapper "github.com/max-grape/test-url-benchmark/wrapper/parser"
	"github.com/stretchr/testify/assert"
)

type parserMock struct {
	data omap.OrderedMap
	err  error
}

func (p *parserMock) Run(searchParam string) (omap.OrderedMap, error) {
	return p.data, p.err
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
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxErr: errors.New("some error"),
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "failed to set processing `parser:foo` key to cache: some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_ParserRunError(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{
			err: errors.New("some error"),
		}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: true,
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_SetError(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: true,
			setErr:        fmt.Errorf("some error"),
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "failed to set `parser:foo` key to cache: some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_LockedSuccess(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{
			data: omap.OrderedMap{{Key: "foo", Val: "bar"}},
		}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: true,
		}),
	)

	actual, err := pcw.Run("foo")

	assert.NoError(t, err)
	assert.Equal(t, omap.OrderedMap{{Key: "foo", Val: "bar"}}, actual)
}

func TestWrapper_Run_GetError(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getErr:        fmt.Errorf("some error"),
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "failed to get `parser:foo` key from cache: some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_GetValueInProgress(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "inprogress",
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, parserCacheWrapper.ErrInProgress.Error())
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_UnmarshalError(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "invalid_json",
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "failed to unmarhsal `parser:foo` key: invalid character 'i' looking for beginning of value")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_DelError(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   "invalid_json",
			delErr:        errors.New("some error"),
		}),
	)

	actual, err := pcw.Run("foo")

	assert.EqualError(t, err, "failed to unmarhsal `parser:foo` key: invalid character 'i' looking for beginning of value; faield to delete cache `parser:foo` key: some error")
	assert.Equal(t, omap.OrderedMap(nil), actual)
}

func TestWrapper_Run_NotLockedSuccess(t *testing.T) {
	pcw := parserCacheWrapper.New(
		parserCacheWrapper.WithParser(&parserMock{}),
		parserCacheWrapper.WithCache(&cacheMock{
			setNxResponse: false,
			getResponse:   `[{"Key":"foo","Val":"bar"}]`,
		}),
	)

	actual, err := pcw.Run("foo")

	assert.NoError(t, err)
	assert.Equal(t, omap.OrderedMap{{Key: "foo", Val: "bar"}}, actual)
}
