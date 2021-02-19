package client

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ChapterClientTestSuite struct {
	suite.Suite
}

func (s *ChapterClientTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestChapterClientTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterClientTestSuite))
}

const titleID = "bleach"

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")

	_ = os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	config.Load()
	defer func() {
		_ = os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
		config.Load()
	}()

	cc := NewChapterClient()
	res, err := cc.GetChapterList(titleID)

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenOriginServerReturns5xxStatusCode() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleID)).
		Reply(http.StatusInternalServerError)

	cc := NewChapterClient()
	res, err := cc.GetChapterList(titleID)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "origin server error: Server is down: returned status code: 500", err.Error())
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenOriginServerReturnsNull() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleID)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader(constants.NullText)))

	cc := NewChapterClient()
	res, err := cc.GetChapterList(titleID)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsError_WhenOriginServerReturnsBrokenJSONResponse() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleID)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader("some error")))

	cc := NewChapterClient()
	res, err := cc.GetChapterList(titleID)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *ChapterClientTestSuite) TestGetChapterList_ReturnsSuccessfulResponse() {
	defer gock.Off()
	gock.New(buildChapterListEndpoint(titleID)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader(`{"komik":[{"hidden_chapter":686,"judul":"Bleach 686 - Death And Strawberry (tamat)","hidden_komik":"bleach","waktu":"2016-08-18 18:59:58"}]}`)))

	cc := NewChapterClient()
	res, err := cc.GetChapterList(titleID)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(res.Chapters) > 0)
}
