package error

import (
	"fmt"
)

type NotFoundError struct {
	S string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find %s", e.S)
}

func NewNotFoundError(s string) *NotFoundError {
	return &NotFoundError{S: s}
}
