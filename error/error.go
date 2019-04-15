package error

import (
	"fmt"
)

type GenericError struct {
	S string
}

func (e *GenericError) Error() string {
	return e.S
}

func NewGenericError() *GenericError {
	return &GenericError{S: "Something went wrong"}
}

type NotFoundError struct {
	S string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find %s", e.S)
}

func NewNotFoundError(s string) *NotFoundError {
	return &NotFoundError{S: s}
}
