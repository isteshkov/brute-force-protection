package common

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func NewUUIDv4() string {
	return uuid.NewV4().String()
}

func CopyStringInterfaceMap(source map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(source))
	for key, value := range source {
		result[key] = value
	}

	return result
}

//nolint:gosec
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func RandString(length int) (result string) {
	var bytes []byte
	for i := 0; i < length; i++ {
		bytes = append(bytes, letters[RandInt(0, len(letters))])
	}

	return string(bytes)
}
