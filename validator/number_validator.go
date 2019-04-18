package validator

import (
	"fmt"
	"regexp"
	"strings"
)

type NumberValidator struct {
	Field string
	Value *string
}

const regexNumber = `\d+\.?\d*`

func (v NumberValidator) Validate() (bool, string) {
	if v.Value == nil {
		return false, fmt.Sprintf("%s cannot be blank", v.Field)
	}

	latLng := strings.TrimSpace(*v.Value)

	if latLng == "" {
		return false, fmt.Sprintf("%s cannot be blank", v.Field)
	}

	regex := regexp.MustCompile(regexNumber)
	if !regex.MatchString(*v.Value) {
		return false, fmt.Sprintf("%s must be a number", v.Field)
	}

	return true, ""
}

func (v NumberValidator) FieldName() string {
	return v.Field
}
