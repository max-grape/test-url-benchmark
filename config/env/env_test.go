package env_test

import (
	"os"
	"testing"

	"github.com/max-grape/test-revo/config/env"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name     string
	env      map[string]string
	input    interface{}
	expected interface{}
	err      string
}

var testTable = []test{
	{
		name:     "input is nil",
		env:      nil,
		input:    nil,
		expected: nil,
		err:      env.ErrInputIsNil.Error(),
	},
	{
		name:     "input is empty non pointer struct",
		env:      nil,
		input:    struct{}{},
		expected: struct{}{},
		err:      env.ErrInputPointerIsZero.Error(),
	},
	{
		name:     "input is empty struct",
		env:      nil,
		input:    &struct{}{},
		expected: &struct{}{},
		err:      "",
	},
	{
		name:     "input is non struct",
		env:      nil,
		input:    "string",
		expected: "string",
		err:      env.ErrNonStructInput.Error(),
	},
	{
		name: "default value, env is not set",
		env:  nil,
		input: &struct {
			Field string `env:"STRING_ENV"`
		}{
			Field: "string default value",
		},
		expected: &struct {
			Field string `env:"STRING_ENV"`
		}{
			Field: "string default value",
		},
		err: "",
	},
	{
		name: "default value, env is set and empty",
		env: map[string]string{
			"STRING_ENV": "",
		},
		input: &struct {
			Field string `env:"STRING_ENV"`
		}{
			Field: "string default value",
		},
		expected: &struct {
			Field string `env:"STRING_ENV"`
		}{
			Field: "",
		},
		err: "",
	},
}

func TestGet(t *testing.T) {
	testGet(t, testTable)
}

func testGet(t *testing.T, table []test) {
	for _, test := range table {
		for key, value := range test.env {
			assert.NoError(t, os.Setenv(key, value))
			defer func(key string) {
				assert.NoError(t, os.Unsetenv(key))
			}(key)
		}

		err := env.Get(test.input)
		if test.err != "" {
			assert.EqualErrorf(t, err, test.err, test.name)
		} else if err != nil {
			assert.Fail(t, "unexpected error: %+v", err)
		}

		assert.Equalf(t, test.expected, test.input, test.name)
	}
}
