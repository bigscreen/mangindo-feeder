package common

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsSuiteSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (s *UtilsTestSuite) TestGetFormattedChapterNumber_WhenChapterHasZeroAfterComma() {
	assert.Equal(s.T(), "100", GetFormattedChapterNumber(float32(100.0)))
}

func (s *UtilsTestSuite) TestGetFormattedChapterNumber_WhenChapterHasNoZeroAfterComma() {
	assert.Equal(s.T(), "100.1", GetFormattedChapterNumber(float32(100.1)))
}
