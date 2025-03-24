package qstruct

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// setValueToPrimitiveField gets a reflect.Value field and a value to parse and set
func setValueToPrimitiveField(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetFloat(parsed)
	default:
		return fmt.Errorf("unsupported primitive kind: %s", field.Kind())
	}

	return nil
}

func setValueToTimeField(field reflect.Value, tag reflect.StructTag, val string) error {
	timeFormat, ok := tag.Lookup("format")
	if !ok {
		timeFormat = time.RFC3339
	}

	t, err := time.Parse(timeFormat, val)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
	}
	field.Set(reflect.ValueOf(t))
	return nil
}
