package qstruct

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

var ErrRequired = errors.New("required field")
var ErrUnexpectedValue = errors.New("unexpected value")
var ErrUnexpectedType = errors.New("unexpected type")

func NewFor[T any](query url.Values) (*T, error) {
	typ := reflect.TypeFor[T]()
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedType, typ.String())
	}

	return hydrate[T](typ, query)
}

func hydrate[T any](typ reflect.Type, query url.Values) (*T, error) {
	result := reflect.New(typ).Interface().(*T)
	reflection := reflect.ValueOf(result).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typField := typ.Field(i)

		name := getFieldName(typField)
		if name == "-" {
			continue
		}

		field := reflection.Field(i)
		if field.Type() != reflect.TypeOf(time.Time{}) {
			if fieldKind := field.Kind(); fieldKind == reflect.Struct || fieldKind == reflect.Array || fieldKind == reflect.Slice || fieldKind == reflect.Map {
				name = fmt.Sprintf("%s[]", name)
			}
		}

		if field.IsValid() && field.CanSet() {
			var wasSet bool
			if values, ok := query[name]; ok {
				if err := setValueToField(field, typField.Tag, values); err != nil {
					return nil, err
				}
				wasSet = true
			} else if val, ok := typField.Tag.Lookup("default"); ok {
				if err := setValueToField(field, typField.Tag, []string{val}); err != nil {
					return nil, err
				}
				wasSet = true
			}

			if val, ok := typField.Tag.Lookup("validate"); ok {
				if slices.Contains(strings.Split(val, ","), "required") && !wasSet {
					return nil, fmt.Errorf("%w: %s", ErrRequired, typField.Name)
				}
			}
		}
	}

	return result, nil
}

func getFieldName(typField reflect.StructField) string {
	if val, ok := typField.Tag.Lookup("query"); ok {
		return val
	}

	return typField.Name
}

func setValueToField(field reflect.Value, tag reflect.StructTag, val []string) error {
	if field.Type() == reflect.TypeOf(time.Time{}) {
		timeFormat, ok := tag.Lookup("format")
		if !ok {
			timeFormat = time.RFC3339
		}

		t, err := time.Parse(timeFormat, val[0])
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.Set(reflect.ValueOf(t))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(val[0])
	case reflect.Bool:
		parsed, err := strconv.ParseBool(val[0])
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(val[0], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(val[0], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val[0], 64)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUnexpectedValue, err)
		}
		field.SetFloat(parsed)
	default:
		return fmt.Errorf("unsupported kind: %s", field.Kind())
	}

	return nil
}
