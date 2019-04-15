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

type ChapterHandlerTestSuite struct {
	suite.Suite
}

func (s *ChapterHandlerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestChapterHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterHandlerTestSuite))
}

func buildChapterParam(titleId string) string {
	param := fmt.Sprintf("?title_id=%s", titleId)
	return constants.GetChaptersApiPath + param
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenQueryParamIsEmpty() {
	req, _ := http.NewRequest("GET", buildChapterParam(""), nil)

	rr := httptest.NewRecorder()

	cs := service.MockChapterService{}

	h := GetChapters(cs)
	h.ServeHTTP(rr, req)

	body := string(rr.Body.Bytes())

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), body, "title_id cannot be blank")
	cs.AssertNotCalled(s.T(), "GetChapters", req.Context(), mock.Anything)
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenUnknownErrorHappens() {
	req, _ := http.NewRequest("GET", buildChapterParam("foo"), nil)

	rr := httptest.NewRecorder()
	err := mErr.NewGenericError()

	cs := service.MockChapterService{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(nil, err)

	h := GetChapters(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenChaptersDoNotExist() {
	req, _ := http.NewRequest("GET", buildChapterParam("foo"), nil)

	rr := httptest.NewRecorder()
	err := mErr.NewNotFoundError("chapter")

	cs := service.MockChapterService{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(nil, err)

	h := GetChapters(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsSuccess_WhenChaptersExist() {
	req, _ := http.NewRequest("GET", buildChapterParam("foo"), nil)

	rr := httptest.NewRecorder()
	cc := contract.Chapter{
		Number:  "54",
		Title:   "Foo",
		TitleId: "foo",
	}
	ccs := []contract.Chapter{cc}
	cr := contract.ChapterResponse{
		Success:  true,
		Chapters: ccs,
	}

	cs := service.MockChapterService{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(&ccs, nil)

	h := GetChapters(cs)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(cr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	cs.AssertExpectations(s.T())
}
