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

type ContentClientTestSuite struct {
	suite.Suite
}

func (s *ContentClientTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestContentClientTestSuite(t *testing.T) {
	suite.Run(t, new(ContentClientTestSuite))
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")

	_ = os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	config.Load()
	defer func() {
		_ = os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
		config.Load()
	}()

	cc := NewContentClient()
	res, err := cc.GetContentList("bleach", 657.0)

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), res)
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsError_WhenOriginServerReturns5xxStatusCode() {
	defer gock.Off()
	gock.New(buildContentListEndpoint("bleach", 657.0)).
		Reply(http.StatusInternalServerError)

	cc := NewContentClient()
	res, err := cc.GetContentList("bleach", 657.0)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "origin server error: Server is down: returned status code: 500", err.Error())
	assert.Nil(s.T(), res)
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsError_WhenOriginServerReturnsNull() {
	defer gock.Off()
	gock.New(buildContentListEndpoint("bleach", 657.0)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader(constants.NullText)))

	cc := NewContentClient()
	res, err := cc.GetContentList("bleach", 657.0)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsError_WhenOriginServerReturnsBrokenJSONResponse() {
	defer gock.Off()
	gock.New(buildContentListEndpoint("bleach", 657.0)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader("some error")))

	cc := NewContentClient()
	res, err := cc.GetContentList("bleach", 657.0)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *ContentClientTestSuite) TestGetContentList_ReturnsSuccessfulResponse() {
	defer gock.Off()
	gock.New(buildContentListEndpoint("bleach", 657.0)).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader(`{"chapter":[{"url":"http://mangacanblog.com/mangas/bleach/657 - thunder god 2/mangacanblogcom_bleach_657_01.jpg","page":1},{"url":"http://mangacanblog.com/mangas/bleach/657 - thunder god 2/mangacanblogcom_bleach_657_02.jpg","page":2},{"url":"http://mangacanblog.com/mangas/bleach/657 - thunder god 2/mangacanblogcom_bleach_657_03.jpg","page":3}]}`)))

	cc := NewContentClient()
	res, err := cc.GetContentList("bleach", 657.0)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(res.Contents) > 0)
}
