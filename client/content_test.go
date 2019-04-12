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

type ContentClientTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ContentClientTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()

	s.ctx = context.Background()
}

func TestContentClientTestSuite(t *testing.T) {
	suite.Run(t, new(ContentClientTestSuite))
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")

	os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	config.Load()

	cc := NewContentClient()
	res, err := cc.GetContentList(s.ctx, "bleach", 657.0)

	os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
	config.Load()

	assert.Contains(s.T(), strings.ToUpper(err.Error()), "TIMEOUT")
	assert.Nil(s.T(), res)
}
