package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	mMock "github.com/bigscreen/mangindo-feeder/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ContentHandlerTestSuite struct {
	suite.Suite
	mr *mux.Router
}

func TestContentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentHandlerTestSuite))
}

func buildContentRequest(titleID, chapter string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", buildContentPath(titleID, chapter), nil)
	rr := httptest.NewRecorder()
	return req, rr
}

func buildContentPath(titleID, chapter string) string {
	tVar := fmt.Sprintf("{%s}", constants.TitleIDKeyParam)
	cVar := fmt.Sprintf("{%s}", constants.ChapterKeyParam)
	replacer := strings.NewReplacer(tVar, titleID, cVar, chapter)
	return replacer.Replace(constants.GetContentsAPIPath)
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
	cs := &mMock.ContentServiceMock{}

	req, rr := buildContentRequest(" ", " ")

	s.mr.HandleFunc(constants.GetContentsAPIPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	body := rr.Body.String()

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "title_id cannot be blank")
	assert.Contains(s.T(), body, "chapter cannot be blank")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenInvalidChapterIsBeingSent() {
	cs := &mMock.ContentServiceMock{}

	req, rr := buildContentRequest("foo", "abc")

	s.mr.HandleFunc(constants.GetContentsAPIPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), rr.Body.String(), "chapter must be a number")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenUnknownErrorHappens() {
	err := mErr.NewGenericError()
	cs := &mMock.ContentServiceMock{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsAPIPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenContentsDoNotExist() {
	err := mErr.NewNotFoundError("chapter")
	cs := &mMock.ContentServiceMock{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsAPIPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsSuccess_WhenContentsExist() {
	cc := contract.Content{
		ImageURL: "http://foo.com/foo.jpg",
	}
	ccs := []contract.Content{cc}
	cr := contract.ContentResponse{
		Success:  true,
		Contents: ccs,
	}
	cs := &mMock.ContentServiceMock{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(&ccs, nil)

	req, rr := buildContentRequest("foo", "123")

	s.mr.HandleFunc(constants.GetContentsAPIPath, GetContents(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(cr)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}
