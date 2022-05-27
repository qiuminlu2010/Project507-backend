package util

import (
	uuid "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	// or error handling
	u2 := uuid.NewV4()
	return u2.String()
}
