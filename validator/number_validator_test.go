package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NumberValidatorTestSuite struct {
	suite.Suite
}

func TestNumberValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(NumberValidatorTestSuite))
}

func (s *NumberValidatorTestSuite) TestValidate_ReturnsFalse_WhenFieldIsMissing() {
	validator := NumberValidator{Field: "foo", Value: nil}
	valid, err := validator.Validate()

	assert.False(s.T(), valid)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "foo cannot be blank", err)
}

func (s *NumberValidatorTestSuite) TestValidate_ReturnsFalse_WhenFieldIsBlank() {
	value := ""
	validator := NumberValidator{Field: "foo", Value: &value}
	valid, err := validator.Validate()

	assert.False(s.T(), valid)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "foo cannot be blank", err)
}

func (s *NumberValidatorTestSuite) TestValidate_ReturnsFalse_WhenInvalidCharactersBeingSent() {
	value := "abc"
	validator := NumberValidator{Field: "foo", Value: &value}
	valid, err := validator.Validate()

	assert.False(s.T(), valid)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "foo must be a number", err)
}

func (s *NumberValidatorTestSuite) TestValidate_ReturnsTrue_WhenFieldIsDecimal() {
	value := "1.23"
	validator := NumberValidator{Field: "foo", Value: &value}
	valid, err := validator.Validate()

	assert.True(s.T(), valid)
	assert.Empty(s.T(), err)
}

func (s *NumberValidatorTestSuite) TestValidate_ReturnsTrue_WhenFieldIsNonDecimal() {
	value := "123"
	validator := NumberValidator{Field: "foo", Value: &value}
	valid, err := validator.Validate()

	assert.True(s.T(), valid)
	assert.Empty(s.T(), err)
}
