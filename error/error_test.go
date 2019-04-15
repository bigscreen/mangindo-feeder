package error

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ErrorTestSuite struct {
	suite.Suite
}

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) TestGetStatusCodeOf_Returns500() {
	err := NewGenericError()
	code := GetStatusCodeOf(err)

	assert.Equal(s.T(), http.StatusInternalServerError, code)
}

func (s *ErrorTestSuite) TestGetStatusCodeOf_Returns404() {
	err := NewNotFoundError("foo")
	code := GetStatusCodeOf(err)

	assert.Equal(s.T(), http.StatusNotFound, code)
}
