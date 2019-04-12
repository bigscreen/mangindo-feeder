package contract

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ContentRequestTestSuite struct {
	suite.Suite
}

func TestContentRequestSuiteSuite(t *testing.T) {
	suite.Run(t, new(ContentRequestTestSuite))
}

func (s *ContentRequestTestSuite) TestNewContentRequest_ReturnsContentRequestWithZeroChapter_WhenInvalidChapterIsBeingSet() {
	req := NewContentRequest("bleach", "1bc")

	assert.Equal(s.T(), "bleach", req.TitleId)
	assert.Equal(s.T(), float32(0.0), req.Chapter)
}

func (s *ContentRequestTestSuite) TestNewContentRequest_ReturnsContentRequestWithZeroLatLng_WhenUnparseableLatLngAreBeingSet() {
	badF := "142524353634252534526262625362625362526727257326573562536253625632563253625362145625362536256325632536" +
		"253621425243536342525345262626253626253625267272573265735625353621425243536342525345262626253626253625267272" +
		"5732657356253625362563256325362536214252435363425253452626262536262536252672725732657356253625362563256325.0"
	req := NewContentRequest("bleach", badF)

	assert.Equal(s.T(), "bleach", req.TitleId)
	assert.Equal(s.T(), float32(0.0), req.Chapter)
}

func (s *ContentRequestTestSuite) TestNewContentRequest_ReturnsValidContentRequest() {
	req := NewContentRequest("bleach", "650")

	assert.Equal(s.T(), "bleach", req.TitleId)
	assert.Equal(s.T(), float32(650.0), req.Chapter)
}
