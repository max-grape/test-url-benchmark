package env_test

import (
	"testing"
)

var testTableString = []test{
	{
		name: "string",
		env: map[string]string{
			"STRING_ENV": "string",
		},
		input: &struct {
			Field   string  `env:"STRING_ENV"`
			Pointer *string `env:"STRING_ENV"`
		}{},
		expected: &struct {
			Field   string  `env:"STRING_ENV"`
			Pointer *string `env:"STRING_ENV"`
		}{
			Field:   "string",
			Pointer: &[]string{"string"}[0],
		},
		err: "",
	},
}

func TestGet_string(t *testing.T) {
	testGet(t, testTableString)
}
