package z9

import (
	"github.com/gofrs/uuid"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// NewUUID wrap
func NewUUID() uuid.UUID {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid
}
