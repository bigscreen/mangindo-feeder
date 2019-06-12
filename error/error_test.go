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

func (s *ErrorTestSuite) TestError_ReturnsGenericError() {
	err := NewGenericError()

	assert.Equal(s.T(), "Something went wrong", err.Error())
}

func (s *ErrorTestSuite) TestError_ReturnsNotFoundError() {
	err := NewNotFoundError("Foo")

	assert.Equal(s.T(), "Could not find Foo", err.Error())
}

func (s *ErrorTestSuite) TestError_ReturnsWorkerError() {
	err := NewWorkerError("Foo")

	assert.Equal(s.T(), "Failed to enqueue job with error: Foo", err.Error())
}

func (s *ErrorTestSuite) TestError_ReturnsValidationError() {
	err := NewValidationError(map[string]string{
		"foo": "foo error",
		"bar": "bar error",
	})

	assert.Contains(s.T(), err.Error(), "foo error")
	assert.Contains(s.T(), err.Error(), "bar error")
}

func (s *ErrorTestSuite) TestGetStatusCodeOf_Returns500() {
	err := NewGenericError()
	code := GetStatusCodeOf(err)

	assert.Equal(s.T(), http.StatusInternalServerError, code)
}

func (s *ErrorTestSuite) TestGetStatusCodeOf_Returns400() {
	err := NewValidationError(map[string]string{})
	code := GetStatusCodeOf(err)

	assert.Equal(s.T(), http.StatusBadRequest, code)
}

func (s *ErrorTestSuite) TestGetStatusCodeOf_Returns404() {
	err := NewNotFoundError("foo")
	code := GetStatusCodeOf(err)

	assert.Equal(s.T(), http.StatusNotFound, code)
}
