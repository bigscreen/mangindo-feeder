package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidatorSuite struct {
	suite.Suite
}

type Payload struct {
	Name string `json:"name"`
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSuite))
}

func (s *ValidatorSuite) TestValidateAll_ReturnsFalse_WhenValidationFails() {
	var validators []Validator
	p := Payload{}
	validators = append(validators, PresenceValidator{Field: "name", Value: &p.Name})
	isValid, _ := ValidateAll(validators)

	assert.False(s.T(), isValid)
}

func (s *ValidatorSuite) TestValidateAll_ReturnsTrue_WhenValidationSucceeds() {
	var validators []Validator
	p := Payload{}
	p.Name = "foo"
	validators = append(validators, PresenceValidator{Field: "name", Value: &p.Name})
	isValid, _ := ValidateAll(validators)

	assert.True(s.T(), isValid)
}
