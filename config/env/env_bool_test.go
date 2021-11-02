package env_test

import (
	"testing"
)

var testTableBool = []test{
	{
		name: "bool",
		env: map[string]string{
			"BOOL_ENV": "true",
		},
		input: &struct {
			Field   bool  `env:"BOOL_ENV"`
			Pointer *bool `env:"BOOL_ENV"`
		}{},
		expected: &struct {
			Field   bool  `env:"BOOL_ENV"`
			Pointer *bool `env:"BOOL_ENV"`
		}{
			Field:   true,
			Pointer: &[]bool{true}[0],
		},
		err: "",
	},
	{
		name: "bool invalid",
		env: map[string]string{
			"BOOL_ENV": "123",
		},
		input: &struct {
			Field bool `env:"BOOL_ENV"`
		}{},
		expected: &struct {
			Field bool `env:"BOOL_ENV"`
		}{},
		err: "field `Field` failed to parse `123` bool: strconv.ParseBool: parsing \"123\": invalid syntax",
	},
	{
		name: "bool pointer (invalid)",
		env: map[string]string{
			"BOOL_ENV": "123",
		},
		input: &struct {
			Pointer *bool `env:"BOOL_ENV"`
		}{},
		expected: &struct {
			Pointer *bool `env:"BOOL_ENV"`
		}{},
		err: "field `Pointer` failed to parse `123` bool: strconv.ParseBool: parsing \"123\": invalid syntax",
	},
}

func TestGet_bool(t *testing.T) {
	testGet(t, testTableBool)
}
