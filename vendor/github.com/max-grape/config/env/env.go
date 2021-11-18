package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	tag = "env"
)

var (
	ErrInputIsNil         = errors.New("input is nil")
	ErrInputPointerIsZero = errors.New("input pointer is zero")
	ErrNonStructInput     = errors.New("non struct input")
)

func Parse(input interface{}) error {
	if input == nil {
		return ErrInputIsNil
	}

	inputPointer := reflect.ValueOf(input)
	if inputPointer.IsZero() {
		return ErrInputPointerIsZero
	}

	inputValue := reflect.Indirect(inputPointer)

	inputType := inputValue.Type()
	if inputType.Kind() != reflect.Struct {
		return ErrNonStructInput
	}

	for i := 0; i < inputValue.NumField(); i++ {
		inputTypeField := inputType.Field(i)

		env, ok := inputTypeField.Tag.Lookup(tag)
		if !ok || env == "" {
			continue
		}

		envValue, ok := os.LookupEnv(env)
		if !ok {
			continue
		}

		inputField := inputValue.Field(i)
		if !inputField.CanSet() {
			return fmt.Errorf("can not set `%s` field", inputTypeField.Name)
		}

		switch inputField.Interface().(type) {

		case string:
			inputField.SetString(envValue)

		case *string:
			inputField.Set(reflect.ValueOf(&envValue))

		case int, int64:
			v, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` int: %w", inputTypeField.Name, envValue, err)
			}
			if inputField.OverflowInt(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			inputField.SetInt(v)

		case *int:
			v, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` int: %w", inputTypeField.Name, envValue, err)
			}
			if !inputField.IsZero() && inputField.Elem().OverflowInt(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			i := int(v)
			inputField.Set(reflect.ValueOf(&i))

		case *int64:
			v, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` int: %w", inputTypeField.Name, envValue, err)
			}
			if !inputField.IsZero() && inputField.Elem().OverflowInt(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			inputField.Set(reflect.ValueOf(&v))

		case uint, uint64:
			v, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` uint: %w", inputTypeField.Name, envValue, err)
			}
			if inputField.OverflowUint(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			inputField.SetUint(v)

		case *uint:
			v, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` uint: %w", inputTypeField.Name, envValue, err)
			}
			if !inputField.IsZero() && inputField.Elem().OverflowUint(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			i := uint(v)
			inputField.Set(reflect.ValueOf(&i))

		case *uint64:
			v, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` uint: %w", inputTypeField.Name, envValue, err)
			}
			if !inputField.IsZero() && inputField.Elem().OverflowUint(v) {
				return fmt.Errorf("field `%s` is oveflown by %d", inputTypeField.Name, v)
			}
			inputField.Set(reflect.ValueOf(&v))

		case bool:
			v, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` bool: %w", inputTypeField.Name, envValue, err)
			}
			inputField.SetBool(v)

		case *bool:
			v, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` bool: %w", inputTypeField.Name, envValue, err)
			}
			inputField.Set(reflect.ValueOf(&v))

		case time.Duration:
			v, err := time.ParseDuration(envValue)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` duration: %w", inputTypeField.Name, envValue, err)
			}
			inputField.Set(reflect.ValueOf(v))

		case *time.Duration:
			v, err := time.ParseDuration(envValue)
			if err != nil {
				return fmt.Errorf("field `%s` failed to parse `%s` duration: %w", inputTypeField.Name, envValue, err)
			}
			inputField.Set(reflect.ValueOf(&v))
		}
	}

	return nil
}
