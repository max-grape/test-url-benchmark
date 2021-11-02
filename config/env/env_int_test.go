package env_test

import (
	"testing"
)

var testTableInt = []test{
	{
		name: "int",
		env: map[string]string{
			"INT_ENV": "123",
		},
		input: &struct {
			Field   int  `env:"INT_ENV"`
			Pointer *int `env:"INT_ENV"`
		}{},
		expected: &struct {
			Field   int  `env:"INT_ENV"`
			Pointer *int `env:"INT_ENV"`
		}{
			Field:   123,
			Pointer: &[]int{123}[0],
		},
		err: "",
	},
	{
		name: "int (invalid)",
		env: map[string]string{
			"INT_ENV": "foo",
		},
		input: &struct {
			Field int `env:"INT_ENV"`
		}{},
		expected: &struct {
			Field int `env:"INT_ENV"`
		}{},
		err: "field `Field` failed to parse `foo` int: strconv.ParseInt: parsing \"foo\": invalid syntax",
	},
	{
		name: "int pointer (invalid)",
		env: map[string]string{
			"INT_ENV": "foo",
		},
		input: &struct {
			Pointer *int `env:"INT_ENV"`
		}{},
		expected: &struct {
			Pointer *int `env:"INT_ENV"`
		}{},
		err: "field `Pointer` failed to parse `foo` int: strconv.ParseInt: parsing \"foo\": invalid syntax",
	},
	{
		name: "int (overflown)",
		env: map[string]string{
			"INT_ENV": "9223372036854775808",
		},
		input: &struct {
			Field int `env:"INT_ENV"`
		}{},
		expected: &struct {
			Field int `env:"INT_ENV"`
		}{},
		err: "field `Field` failed to parse `9223372036854775808` int: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
	},
	{
		name: "int pointer (overflown)",
		env: map[string]string{
			"INT_ENV": "9223372036854775808",
		},
		input: &struct {
			Pointer *int `env:"INT_ENV"`
		}{},
		expected: &struct {
			Pointer *int `env:"INT_ENV"`
		}{},
		err: "field `Pointer` failed to parse `9223372036854775808` int: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
	},
}

func TestGet_int(t *testing.T) {
	testGet(t, testTableInt)
}
