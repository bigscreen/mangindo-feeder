package client

import (
	"context"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type ChapterClientTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ChapterClientTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()

	s.ctx = context.Background()
}

func TestChapterClientTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterClientTestSuite))
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")
	os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	defer os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
	config.Load()

	cc := NewChapterClient()
	res, err := cc.GetChapterList(s.ctx, "bleach")

	assert.Contains(s.T(), strings.ToUpper(err.Error()), "TIMEOUT")
	assert.Nil(s.T(), res)
}
