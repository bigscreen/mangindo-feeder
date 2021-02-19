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

type ChapterHandlerTestSuite struct {
	suite.Suite
	mr *mux.Router
}

func TestChapterHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterHandlerTestSuite))
}

func buildChapterRequest(titleID string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest("GET", buildChapterPath(titleID), nil)
	rr := httptest.NewRecorder()
	return req, rr
}

func buildChapterPath(titleID string) string {
	pVar := fmt.Sprintf("{%s}", constants.TitleIDKeyParam)
	return strings.Replace(constants.GetChaptersAPIPath, pVar, titleID, -1)
}

func (s *ChapterHandlerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ChapterHandlerTestSuite) SetupTest() {
	s.mr = mux.NewRouter()
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenQueryParamIsBlank() {
	cs := &mMock.ChapterServiceMock{}

	req, rr := buildChapterRequest(" ")

	s.mr.HandleFunc(constants.GetChaptersAPIPath, GetChapters(cs))
	s.mr.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusBadRequest, rr.Code)
	assert.Contains(s.T(), rr.Body.String(), "title_id cannot be blank")
	cs.AssertNotCalled(s.T(), "GetChapters", req.Context(), mock.Anything)
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenUnknownErrorHappens() {
	err := mErr.NewGenericError()
	cs := &mMock.ChapterServiceMock{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(nil, err)

	req, rr := buildChapterRequest("foo")

	s.mr.HandleFunc(constants.GetChaptersAPIPath, GetChapters(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsError_WhenChaptersDoNotExist() {
	err := mErr.NewNotFoundError("chapter")
	cs := &mMock.ChapterServiceMock{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(nil, err)

	req, rr := buildChapterRequest("foo")

	s.mr.HandleFunc(constants.GetChaptersAPIPath, GetChapters(cs))
	s.mr.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}

func (s *ChapterHandlerTestSuite) TestGetChapters_ReturnsSuccess_WhenChaptersExist() {
	cc := contract.Chapter{
		Number:  "54",
		Title:   "Foo",
		TitleID: "foo",
	}
	ccs := []contract.Chapter{cc}
	cs := &mMock.ChapterServiceMock{}
	cs.On("GetChapters", contract.NewChapterRequest("foo")).Return(&ccs, nil)

	req, rr := buildChapterRequest("foo")

	s.mr.HandleFunc(constants.GetChaptersAPIPath, GetChapters(cs))
	s.mr.ServeHTTP(rr, req)

	cr := contract.ChapterResponse{
		Success:  true,
		Chapters: ccs,
	}
	res, _ := json.Marshal(cr)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), strings.TrimSuffix(rr.Body.String(), "\n"))
	cs.AssertExpectations(s.T())
}
