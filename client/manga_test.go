package client

import (
	"github.com/ad2games/vcr-go"
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

type MangaClientTestSuite struct {
	suite.Suite
}

func (s *MangaClientTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestMangaClientTestSuite(t *testing.T) {
	suite.Run(t, new(MangaClientTestSuite))
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsError_WhenCallTimesOut() {
	ht := os.Getenv("HYSTRIX_TIMEOUT_MS")

	_ = os.Setenv("HYSTRIX_TIMEOUT_MS", "1")
	config.Load()
	defer func() {
		_ = os.Setenv("HYSTRIX_TIMEOUT_MS", ht)
		config.Load()
	}()

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.Contains(s.T(), strings.ToUpper(err.Error()), "TIMEOUT")
	assert.Nil(s.T(), res)
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsError_WhenOriginServerReturns5xxStatusCode() {
	defer gock.Off()
	gock.New(buildMangaListEndpoint()).
		Reply(http.StatusInternalServerError)

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "origin server error: Server is down: returned status code: 500", err.Error())
	assert.Nil(s.T(), res)
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsError_WhenOriginServerReturnsBrokenJSONResponse() {
	defer gock.Off()
	gock.New(buildMangaListEndpoint()).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader("some error")))

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsError_WhenOriginServerReturnsNull() {
	defer gock.Off()
	gock.New(buildMangaListEndpoint()).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader("null")))

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsSuccessfulResponse() {
	vcr.Start("get_manga_list_valid_response", nil)
	defer vcr.Stop()

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(res.Mangas) > 0)
}
