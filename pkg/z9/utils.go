package z9

import (
	"github.com/gofrs/uuid"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type UUID = uuid.UUID

// NewUUID wrap
func NewUUID() UUID {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid
}
