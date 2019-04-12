package client

import (
	"context"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
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

const titleId = "bleach"

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")
	os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	defer os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
	config.Load()

	cc := NewChapterClient()
	res, err := cc.GetChapterList(s.ctx, titleId)

	assert.Contains(s.T(), strings.ToUpper(err.Error()), "TIMEOUT")
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenOriginServerReturns5xxStatusCode() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleId)).
		Reply(http.StatusInternalServerError)

	cc := NewChapterClient()
	res, err := cc.GetChapterList(s.ctx, titleId)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "origin server error: Server is down: returned status code: 500", err.Error())
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenOriginServerReturnsBrokenJSONResponse() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleId)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader("some error")))

	cc := NewChapterClient()
	res, err := cc.GetChapterList(s.ctx, titleId)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}
