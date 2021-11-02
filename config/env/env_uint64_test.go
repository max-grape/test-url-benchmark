package env_test

import (
	"testing"
)

var testTableUint64 = []test{
	{
		name: "uint64",
		env: map[string]string{
			"UINT64_ENV": "123",
		},
		input: &struct {
			Field   uint64  `env:"UINT64_ENV"`
			Pointer *uint64 `env:"UINT64_ENV"`
		}{},
		expected: &struct {
			Field   uint64  `env:"UINT64_ENV"`
			Pointer *uint64 `env:"UINT64_ENV"`
		}{
			Field:   123,
			Pointer: &[]uint64{123}[0],
		},
		err: "",
	},
	{
		name: "uint64 (invalid)",
		env: map[string]string{
			"UINT64_ENV": "foo",
		},
		input: &struct {
			Field uint64 `env:"UINT64_ENV"`
		}{},
		expected: &struct {
			Field uint64 `env:"UINT64_ENV"`
		}{},
		err: "field `Field` failed to parse `foo` uint: strconv.ParseUint: parsing \"foo\": invalid syntax",
	},
	{
		name: "uint64 pointer (invalid)",
		env: map[string]string{
			"UINT64_ENV": "foo",
		},
		input: &struct {
			Pointer *uint64 `env:"UINT64_ENV"`
		}{},
		expected: &struct {
			Pointer *uint64 `env:"UINT64_ENV"`
		}{},
		err: "field `Pointer` failed to parse `foo` uint: strconv.ParseUint: parsing \"foo\": invalid syntax",
	},
	{
		name: "uint64 (overflown)",
		env: map[string]string{
			"UINT64_ENV": "18446744073709551616",
		},
		input: &struct {
			Field uint64 `env:"UINT64_ENV"`
		}{},
		expected: &struct {
			Field uint64 `env:"UINT64_ENV"`
		}{},
		err: "field `Field` failed to parse `18446744073709551616` uint: strconv.ParseUint: parsing \"18446744073709551616\": value out of range",
	},
	{
		name: "uint64 pointer (overflown)",
		env: map[string]string{
			"UINT64_ENV": "18446744073709551616",
		},
		input: &struct {
			Pointer *uint64 `env:"UINT64_ENV"`
		}{},
		expected: &struct {
			Pointer *uint64 `env:"UINT64_ENV"`
		}{},
		err: "field `Pointer` failed to parse `18446744073709551616` uint: strconv.ParseUint: parsing \"18446744073709551616\": value out of range",
	},
}

func TestGet_uint64(t *testing.T) {
	testGet(t, testTableUint64)
}
