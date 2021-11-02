package env_test

import (
	"testing"
	"time"
)

var testTableDuration = []test{
	{
		name: "duration",
		env: map[string]string{
			"DURATION_ENV": "123s",
		},
		input: &struct {
			Field   time.Duration  `env:"DURATION_ENV"`
			Pointer *time.Duration `env:"DURATION_ENV"`
		}{},
		expected: &struct {
			Field   time.Duration  `env:"DURATION_ENV"`
			Pointer *time.Duration `env:"DURATION_ENV"`
		}{
			Field:   time.Second * 123,
			Pointer: &[]time.Duration{time.Second * 123}[0],
		},
		err: "",
	},
	{
		name: "duration (invalid)",
		env: map[string]string{
			"DURATION_ENV": "foo",
		},
		input: &struct {
			Field time.Duration `env:"DURATION_ENV"`
		}{},
		expected: &struct {
			Field time.Duration `env:"DURATION_ENV"`
		}{},
		err: "field `Field` failed to parse `foo` duration: time: invalid duration \"foo\"",
	},
	{
		name: "duration pointer (invalid)",
		env: map[string]string{
			"DURATION_ENV": "foo",
		},
		input: &struct {
			Pointer *time.Duration `env:"DURATION_ENV"`
		}{},
		expected: &struct {
			Pointer *time.Duration `env:"DURATION_ENV"`
		}{},
		err: "field `Pointer` failed to parse `foo` duration: time: invalid duration \"foo\"",
	},
}

func TestGet_duration(t *testing.T) {
	testGet(t, testTableDuration)
}
