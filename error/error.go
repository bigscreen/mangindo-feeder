package error

import (
	"fmt"
	"net/http"
	"reflect"
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

func GetStatusCodeOf(objectPtr interface{}) int {
	if isErrorInstanceOf(objectPtr, (*NotFoundError)(nil)) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func isErrorInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}
