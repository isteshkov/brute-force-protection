package common

import (
	uuid "github.com/satori/go.uuid"
)

func NewUUIDv4() string {
	return uuid.NewV4().String()
}

func KeysFromStringBoolMap(data map[string]bool) []string {
	result := make([]string, 0, len(data))
	for key := range data {
		result = append(result, key)
	}

	return result
}

func KeysFromStringMap(data map[string]string) []string {
	result := make([]string, 0, len(data))
	for key := range data {
		result = append(result, key)
	}

	return result
}

func ValuesFromStringMap(data map[string]string) []string {
	result := make([]string, 0, len(data))
	for _, value := range data {
		result = append(result, value)
	}

	return result
}

func CopyStringInterfaceMap(source map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(source))
	for key, value := range source {
		result[key] = value
	}

	return result
}

func RemoveFromStringsSlice(s []string, elem string) []string {
	for i, v := range s {
		if v == elem {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func NonexistentStringsInMapStringStruct(source []string, searchMap map[string]struct{}) (nonexistent []string) {
	for i := range source {
		if _, ok := searchMap[source[i]]; !ok {
			nonexistent = append(nonexistent, source[i])
		}
	}
	return
}
