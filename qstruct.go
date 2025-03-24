package qstruct

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"slices"
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

	result := reflect.New(typ)
	for i := 0; i < typ.NumField(); i++ {
		typField := typ.Field(i)
		field := result.Elem().Field(i)
		if field.IsValid() && field.CanSet() {
			if queryPath := getQueryPath(typField); queryPath != "-" {
				if err := hydrate(query, field, typField, queryPath, 0); err != nil {
					return nil, err
				}
			}
		}
	}

	return result.Interface().(*T), nil
}

func hydrate(query url.Values, field reflect.Value, typField reflect.StructField, queryPath string, queryValueIndex int) error {
	wasSet, err := setFieldValue(query, field, typField, queryPath, queryValueIndex)
	if err != nil {
		return err
	}

	// Handle validations
	if val, ok := typField.Tag.Lookup("validate"); ok {
		if slices.Contains(strings.Split(val, ","), "required") && !wasSet {
			return fmt.Errorf("%w: %s", ErrRequired, typField.Name)
		}
	}

	return nil
}

func setFieldValue(query url.Values, field reflect.Value, typField reflect.StructField, queryPath string, queryValueIndex int) (bool, error) {
	if field.Kind() == reflect.Slice {
		if err := setValueToSliceField(query, field, typField, queryPath); err != nil {
			return false, err
		}

		return true, nil
	}

	if values, ok := query[queryPath]; ok {
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if err := setValueToTimeField(field, typField.Tag, values[queryValueIndex]); err != nil {
				return false, err
			}
			return true, nil
		}
		if err := setValueToPrimitiveField(field, values[queryValueIndex]); err != nil {
			return false, err
		}
		return true, nil
	}

	// Handle default values
	if val, ok := typField.Tag.Lookup("default"); ok {
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if err := setValueToTimeField(field, typField.Tag, val); err != nil {
				return false, err
			}
			return true, nil
		}
		if err := setValueToPrimitiveField(field, val); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func getQueryPath(typField reflect.StructField) string {
	if val, ok := typField.Tag.Lookup("query"); ok {
		return val
	}

	return typField.Name
}
