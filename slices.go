package qstruct

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrSliceRegexMatchError = errors.New("regex matched unexpected values")

func setValueToSliceField(query url.Values, field reflect.Value, typField reflect.StructField, queryPath string, queryValueIndex int) error {
	return processIndexedSlice(query, field, typField, queryPath, queryValueIndex)
}

func processIndexedSlice(query url.Values, field reflect.Value, typField reflect.StructField, queryPath string, queryValueIndex int) error {
	type matchStruct struct {
		key             string
		targetQueryPath string
		exact           bool
		parsedKey       int
	}

	highestKey := -1
	appends := 0
	matches := make([]matchStruct, 0)
	regex := regexp.MustCompile(fmt.Sprintf("^%s\\[(\\d*)\\]", regexp.QuoteMeta(queryPath)))
	for queryParam, values := range query {
		if found := regex.FindStringSubmatch(queryParam); found != nil {
			if len(found) != 2 {
				return fmt.Errorf("%w: %s", ErrSliceRegexMatchError, queryParam)
			}

			targetQueryPath := fmt.Sprintf("%s[%s]", queryPath, found[1])
			exact := queryParam == targetQueryPath
			parsedKey := getIndexFromMatch(found[1])
			matches = append(matches, matchStruct{
				key:             found[1],
				targetQueryPath: targetQueryPath,
				exact:           exact,
				parsedKey:       parsedKey,
			})

			if found[1] == "" {
				if strings.HasSuffix(targetQueryPath, "[]") && !strings.HasSuffix(queryPath, "[]") {
					appends += len(values)
				} else {
					appends++
				}
			} else if highestKey < parsedKey {
				highestKey = parsedKey
			}
		}
	}

	reflection := reflect.New(field.Type())
	defer field.Set(reflection.Elem())
	size := highestKey + appends + 1
	slice := reflect.MakeSlice(field.Type(), size, size)
	defer reflection.Elem().Set(slice)

	startFrom := 0
	if highestKey >= 0 {
		startFrom = highestKey + 1
	}
	for _, match := range matches {
		if match.key == "" { // Case unindexed
			if strings.HasSuffix(match.targetQueryPath, "[]") && !strings.HasSuffix(queryPath, "[]") { // case nested
				for i := range appends {
					if err := hydrate(query, slice.Index(startFrom+i), typField, match.targetQueryPath, i); err != nil {
						return err
					}
				}
				continue
			}

			if err := hydrate(query, slice.Index(startFrom+0), typField, match.targetQueryPath, queryValueIndex); err != nil {
				return err
			}
			continue
		}

		if err := hydrate(query, slice.Index(match.parsedKey), typField, match.targetQueryPath, 0); err != nil {
			return err
		}
	}
	return nil
}

func getIndexFromMatch(found string) int {
	index, err := strconv.Atoi(found)
	if err != nil {
		return 0
	}

	return index
}
