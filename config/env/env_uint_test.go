package env_test

import (
	"testing"
)

var testTableUint = []test{
	{
		name: "uint",
		env: map[string]string{
			"UINT_ENV": "123",
		},
		input: &struct {
			Field   uint  `env:"UINT_ENV"`
			Pointer *uint `env:"UINT_ENV"`
		}{},
		expected: &struct {
			Field   uint  `env:"UINT_ENV"`
			Pointer *uint `env:"UINT_ENV"`
		}{
			Field:   123,
			Pointer: &[]uint{123}[0],
		},
		err: "",
	},
	{
		name: "uint (invalid)",
		env: map[string]string{
			"UINT_ENV": "foo",
		},
		input: &struct {
			Field uint `env:"UINT_ENV"`
		}{},
		expected: &struct {
			Field uint `env:"UINT_ENV"`
		}{},
		err: "field `Field` failed to parse `foo` uint: strconv.ParseUint: parsing \"foo\": invalid syntax",
	},
	{
		name: "uint pointer (invalid)",
		env: map[string]string{
			"UINT_ENV": "foo",
		},
		input: &struct {
			Pointer *uint `env:"UINT_ENV"`
		}{},
		expected: &struct {
			Pointer *uint `env:"UINT_ENV"`
		}{},
		err: "field `Pointer` failed to parse `foo` uint: strconv.ParseUint: parsing \"foo\": invalid syntax",
	},
	{
		name: "uint (overflown)",
		env: map[string]string{
			"UINT_ENV": "18446744073709551616",
		},
		input: &struct {
			Field uint `env:"UINT_ENV"`
		}{},
		expected: &struct {
			Field uint `env:"UINT_ENV"`
		}{},
		err: "field `Field` failed to parse `18446744073709551616` uint: strconv.ParseUint: parsing \"18446744073709551616\": value out of range",
	},
	{
		name: "uint pointer (overflown)",
		env: map[string]string{
			"UINT_ENV": "18446744073709551616",
		},
		input: &struct {
			Pointer *uint `env:"UINT_ENV"`
		}{},
		expected: &struct {
			Pointer *uint `env:"UINT_ENV"`
		}{},
		err: "field `Pointer` failed to parse `18446744073709551616` uint: strconv.ParseUint: parsing \"18446744073709551616\": value out of range",
	},
}

func TestGet_uint(t *testing.T) {
	testGet(t, testTableUint)
}
