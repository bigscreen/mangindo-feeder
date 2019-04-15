package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PresenceValidatorSuite struct {
	suite.Suite
}

func TestPresenceValidatorSuite(t *testing.T) {
	suite.Run(t, new(PresenceValidatorSuite))
}

func (s *PresenceValidatorSuite) TestValidatePresence_ReturnsFalse_WhenFieldIsMissing() {
	presenceValidator := PresenceValidator{Field: "name", Value: nil}
	valid, err := presenceValidator.Validate()

	assert.False(s.T(), valid)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "name cannot be blank", err)
}

func (s *PresenceValidatorSuite) TestValidatePresence_ReturnsFalse_WhenFieldIsBlank() {
	value := ""
	presenceValidator := PresenceValidator{Field: "name", Value: &value}
	valid, err := presenceValidator.Validate()

	assert.False(s.T(), valid)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "name cannot be blank", err)
}

func (s *PresenceValidatorSuite) TestValidatePresence_ReturnsTrue_WhenFieldIsPresent() {
	value := "some-value"
	presenceValidator := PresenceValidator{Field: "name", Value: &value}
	valid, _ := presenceValidator.Validate()

	assert.True(s.T(), valid)
}
