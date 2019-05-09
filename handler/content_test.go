package handler

import (
	"encoding/json"
	"fmt"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type ContentHandlerTestSuite struct {
	suite.Suite
	mr *mux.Router
}

func TestContentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentHandlerTestSuite))
}

func buildContentRequest(titleId, chapter string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", buildContentPath(titleId, chapter), nil)
	rr := httptest.NewRecorder()
	return req, rr
}

func buildContentPath(titleId, chapter string) string {
	tVar := fmt.Sprintf("{%s}", constants.TitleIdKeyParam)
	cVar := fmt.Sprintf("{%s}", constants.ChapterKeyParam)
	replacer := strings.NewReplacer(tVar, titleId, cVar, chapter,)
	return replacer.Replace(constants.GetContentsApiPath)
}

func (s *ContentHandlerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ContentHandlerTestSuite) SetupTest() {
	s.mr = mux.NewRouter()
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenParamsAreBlank() {
	cs := service.MockContentService{}

	req, rr := buildContentRequest(" ", " ")

	s.mr.HandleFunc(constants.GetContentsApiPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	body := string(rr.Body.Bytes())

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "title_id cannot be blank")
	assert.Contains(s.T(), body, "chapter cannot be blank")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenInvalidChapterIsBeingSent() {
	cs := service.MockContentService{}

	req, rr := buildContentRequest("foo", "abc")

	s.mr.HandleFunc(constants.GetContentsApiPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	body := string(rr.Body.Bytes())

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "chapter must be a number")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenUnknownErrorHappens() {
	err := mErr.NewGenericError()
	cs := service.MockContentService{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsApiPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenContentsDoNotExist() {
	err := mErr.NewNotFoundError("chapter")
	cs := service.MockContentService{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsApiPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsSuccess_WhenContentsExist() {
	cc := contract.Content{
		ImageURL:  "http://foo.com/foo.jpg",
	}
	ccs := []contract.Content{cc}
	cr := contract.ContentResponse{
		Success:  true,
		Contents: ccs,
	}
	cs := service.MockContentService{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(&ccs, nil)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsApiPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(cr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}
