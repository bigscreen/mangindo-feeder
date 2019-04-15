package validator

import (
	"fmt"
	"strings"
)

type PresenceValidator struct {
	Field string
	Value *string
}

func (p PresenceValidator) Validate() (bool, string) {
	if p.Value == nil {
		return false, fmt.Sprintf("%s cannot be blank", p.Field)
	}

	stringTrimmed := strings.TrimSpace(*p.Value)
	if stringTrimmed == "" {
		return false, fmt.Sprintf("%s cannot be blank", p.Field)
	}

	return true, ""
}

func (p PresenceValidator) FieldName() string {
	return p.Field
}
