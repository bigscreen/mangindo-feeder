package client

import (
	"encoding/json"
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

	assert.NotNil(s.T(), err)
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
		Body(ioutil.NopCloser(strings.NewReader(constants.NullText)))

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), constants.InvalidJSONResponseError, err.Error())
	assert.Nil(s.T(), res)
}

func (s *MangaClientTestSuite) TestGetMangaList_ReturnsSuccessfulResponse() {
	defer gock.Off()
	gock.New(buildMangaListEndpoint()).
		Reply(http.StatusOK).
		Body(ioutil.NopCloser(strings.NewReader(`{"komik":[{"id":"1","judul":"Boku No Hero Academia","hidden_komik":"boku_no_hero_academia","icon_komik":"http://www.mangacanblog.com/official/img/boku_no_hero_academia.jpg","hiddenNewChapter":"224","lastModified":"2019-04-12 15:28:03","genre":"Action, Adventure, Comedy, Shounen, School Life, Sci-Fi, Supernatural","nama_lain":"Boku No Hero Academia","pengarang":"Horikoshi Kouhei","status":"OnGoing","published":"2014","summary":"Cerita ditetapkan di hari modern, kecuali orang-orang dengan kekuatan spesial di seluruh dunia. Anak laki-laki bernama Izuku Modoriya tidak memiliki kekuatan, tapi dia masih bermimpi, Penasaran? simak kisahnya hanya di mangacanblog.com."}]}`)))

	mc := NewMangaClient()
	res, err := mc.GetMangaList()

	a, _ := json.Marshal(res)
	println(string(a))

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(res.Mangas) > 0)
}
