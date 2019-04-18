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
}

func (s *ContentHandlerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestContentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentHandlerTestSuite))
}

func buildContentParam(titleId, chapter string) string {
	param := fmt.Sprintf("?title_id=%s&chapter=%s", titleId, chapter)
	return constants.GetContentsApiPath + param
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenParamsAreEmpty() {
	req, _ := http.NewRequest("GET", buildContentParam("", ""), nil)

	rr := httptest.NewRecorder()

	cs := service.MockContentService{}

	h := GetContents(cs)
	h.ServeHTTP(rr, req)

	body := string(rr.Body.Bytes())

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "title_id cannot be blank")
	assert.Contains(s.T(), body, "chapter cannot be blank")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenInvalidChapterIsBeingSent() {
	req, _ := http.NewRequest("GET", buildContentParam("foo", "abc"), nil)

	rr := httptest.NewRecorder()

	cs := service.MockContentService{}

	h := GetContents(cs)
	h.ServeHTTP(rr, req)

	body := string(rr.Body.Bytes())

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "chapter must be a number")
	cs.AssertNotCalled(s.T(), "GetContents", req.Context(), mock.Anything)
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenUnknownErrorHappens() {
	req, _ := http.NewRequest("GET", buildContentParam("foo", "123"), nil)

	rr := httptest.NewRecorder()
	err := mErr.NewGenericError()

	cs := service.MockContentService{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	h := GetContents(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsError_WhenContentsDoNotExist() {
	req, _ := http.NewRequest("GET", buildContentParam("foo", "123"), nil)

	rr := httptest.NewRecorder()
	err := mErr.NewNotFoundError("chapter")

	cs := service.MockContentService{}
	cs.On("GetContents", contract.NewContentRequest("foo", "123")).Return(nil, err)

	h := GetContents(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ContentHandlerTestSuite) TestGetContents_ReturnsSuccess_WhenContentsExist() {
	req, _ := http.NewRequest("GET", buildContentParam("foo", "123"), nil)

	rr := httptest.NewRecorder()
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

	h := GetContents(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(cr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}
