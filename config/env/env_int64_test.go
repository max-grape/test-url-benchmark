package env_test

import (
	"testing"
)

var testTableInt64 = []test{
	{
		name: "int64",
		env: map[string]string{
			"INT64_ENV": "123",
		},
		input: &struct {
			Field   int64  `env:"INT64_ENV"`
			Pointer *int64 `env:"INT64_ENV"`
		}{},
		expected: &struct {
			Field   int64  `env:"INT64_ENV"`
			Pointer *int64 `env:"INT64_ENV"`
		}{
			Field:   123,
			Pointer: &[]int64{123}[0],
		},
		err: "",
	},
	{
		name: "int64 (invalid)",
		env: map[string]string{
			"INT64_ENV": "foo",
		},
		input: &struct {
			Field int `env:"INT64_ENV"`
		}{},
		expected: &struct {
			Field int `env:"INT64_ENV"`
		}{},
		err: "field `Field` failed to parse `foo` int: strconv.ParseInt: parsing \"foo\": invalid syntax",
	},
	{
		name: "int64 pointer (invalid)",
		env: map[string]string{
			"INT64_ENV": "foo",
		},
		input: &struct {
			Pointer *int `env:"INT64_ENV"`
		}{},
		expected: &struct {
			Pointer *int `env:"INT64_ENV"`
		}{},
		err: "field `Pointer` failed to parse `foo` int: strconv.ParseInt: parsing \"foo\": invalid syntax",
	},
	{
		name: "int64 (overflown)",
		env: map[string]string{
			"INT64_ENV": "9223372036854775808",
		},
		input: &struct {
			Field int64 `env:"INT64_ENV"`
		}{},
		expected: &struct {
			Field int64 `env:"INT64_ENV"`
		}{},
		err: "field `Field` failed to parse `9223372036854775808` int: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
	},
	{
		name: "int64 pointer (overflown)",
		env: map[string]string{
			"INT64_ENV": "9223372036854775808",
		},
		input: &struct {
			Pointer *int64 `env:"INT64_ENV"`
		}{},
		expected: &struct {
			Pointer *int64 `env:"INT64_ENV"`
		}{},
		err: "field `Pointer` failed to parse `9223372036854775808` int: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
	},
}

func TestGet_int64(t *testing.T) {
	testGet(t, testTableInt64)
}
