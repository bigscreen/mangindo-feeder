package error

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

type ValidationError struct {
	validationErrors map[string]string
}

func (e *ValidationError) Error() string {
	var errorMsgs []string
	for _, errMsg := range e.validationErrors {
		errorMsgs = append(errorMsgs, errMsg)
	}
	return strings.Join(errorMsgs, ", ")
}

func NewValidationError(validationErrors map[string]string) *ValidationError {
	return &ValidationError{validationErrors: validationErrors}
}

func GetStatusCodeOf(objectPtr interface{}) int {
	if isErrorInstanceOf(objectPtr, (*NotFoundError)(nil)) {
		return http.StatusNotFound
	} else if isErrorInstanceOf(objectPtr, (*ValidationError)(nil)) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func isErrorInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}
