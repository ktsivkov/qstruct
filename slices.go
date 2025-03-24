package qstruct

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

var ErrSliceRegexMatchError = errors.New("regex matched unexpected values")

func setValueToSliceField(query url.Values, field reflect.Value, typField reflect.StructField, queryPath string) error {
	matchesResult, err := getSliceMatches(query, queryPath)
	if err != nil {
		return err
	}

	reflection := reflect.New(field.Type())
	defer field.Set(reflection.Elem())

	startFrom := matchesResult.highestKey + 1
	size := startFrom + matchesResult.appends
	slice := reflect.MakeSlice(field.Type(), size, size)
	defer reflection.Elem().Set(slice)

	for key, matches := range matchesResult.matches {
		targetQueryParam := fmt.Sprintf("%s[%s]", queryPath, key)
		if key == "" { // Handle unindexed
			for queryParam, match := range matches {
				for _, queryParameter := range match {
					if err := hydrate(url.Values{queryParam: []string{queryParameter}}, slice.Index(startFrom), typField, targetQueryParam, 0); err != nil {
						return err
					}
					startFrom++
				}
			}
			continue
		}

		// Handle indexed
		if err := hydrate(matches, slice.Index(getIndexFromMatch(key)), typField, fmt.Sprintf("%s[%s]", queryPath, key), 0); err != nil {
			return err
		}
	}
	return nil
}

type sliceMatches struct {
	matches    map[string]url.Values
	highestKey int
	appends    int
}

func getSliceMatches(query url.Values, queryPath string) (*sliceMatches, error) {
	regex := regexp.MustCompile(fmt.Sprintf("^%s\\[(\\d*)\\]", regexp.QuoteMeta(queryPath)))
	matchesMap := make(map[string]url.Values)
	highestKey := -1
	appends := 0
	for queryParam, values := range query {
		if found := regex.FindStringSubmatch(queryParam); found != nil {
			if len(found) != 2 {
				return nil, fmt.Errorf("%w: %s", ErrSliceRegexMatchError, queryParam)
			}

			if parsedKey := getIndexFromMatch(found[1]); highestKey < parsedKey {
				highestKey = parsedKey
			}

			if found[1] == "" {
				appends += len(values)
			}

			if prev, ok := matchesMap[found[1]]; ok {
				prev[queryParam] = values
				matchesMap[found[1]] = prev
				continue
			}

			matchesMap[found[1]] = url.Values{
				queryParam: values,
			}
		}
	}
	return &sliceMatches{
		matches:    matchesMap,
		highestKey: highestKey,
		appends:    appends,
	}, nil
}

func getIndexFromMatch(found string) int {
	index, err := strconv.Atoi(found)
	if err != nil {
		return -1
	}

	return index
}
