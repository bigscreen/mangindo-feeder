package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ChapterServiceTestSuite struct {
	suite.Suite
	cca cache.ChapterCache
	cc  *mock.ChapterClientMock
	ws  *mock.WorkerServiceMock
}

func TestChapterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterServiceTestSuite))
}

func (s *ChapterServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ChapterServiceTestSuite) SetupTest() {
	s.cca = cache.NewChapterCache()
	s.cc = &mock.ChapterClientMock{}
	s.ws = &mock.WorkerServiceMock{}
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	ccm := manager.NewChapterCacheManager(s.cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	s.cc.On("GetChapterList", req.TitleID).Return(nil, errors.New("some error"))

	cs := NewChapterService(s.cc, ccm, s.ws)
	cl, err := cs.GetChapters(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	s.cc.AssertExpectations(s.T())
	s.ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleID)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheHitsAndChapterListIsEmpty() {
	ccm := manager.NewChapterCacheManager(s.cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{}}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID)
	}()

	cs := NewChapterService(s.cc, ccm, s.ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), chapters)
	assert.Equal(s.T(), mErr.NewNotFoundError("chapter").Error(), err.Error())

	s.ws.AssertNotCalled(s.T(), "GetChapterList", req.TitleID)
	s.ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleID)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheMissesAndChapterListIsEmpty() {
	ccm := manager.NewChapterCacheManager(s.cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{}}

	s.cc.On("GetChapterList", req.TitleID).Return(&cr, nil)
	s.ws.On("SetChapterCache", req.TitleID).Return(nil)

	cs := NewChapterService(s.cc, ccm, s.ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), chapters)
	assert.Equal(s.T(), mErr.NewNotFoundError("chapter").Error(), err.Error())

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsSuccess_WhenCacheHitsAndChapterListIsNotEmpty() {
	ccm := manager.NewChapterCacheManager(s.cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	dc := domain.Chapter{
		Number:       650.0,
		Title:        "Bleach",
		TitleID:      "bleach",
		ModifiedDate: "2016-08-18 18:59:58",
	}
	cr := &domain.ChapterListResponse{Chapters: []domain.Chapter{dc}}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID)
	}()

	cs := NewChapterService(s.cc, ccm, s.ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*chapters) > 0)
	assert.Equal(s.T(), "650", (*chapters)[0].Number)
	assert.Equal(s.T(), dc.Title, (*chapters)[0].Title)
	assert.Equal(s.T(), dc.TitleID, (*chapters)[0].TitleID)

	s.ws.AssertNotCalled(s.T(), "GetChapterList", req.TitleID)
	s.ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleID)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsSuccess_WhenCacheMissesAndChapterListIsNotEmpty() {
	ccm := manager.NewChapterCacheManager(s.cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	dc := domain.Chapter{
		Number:       650.0,
		Title:        "Bleach",
		TitleID:      "bleach",
		ModifiedDate: "2016-08-18 18:59:58",
	}
	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{dc}}

	s.cc.On("GetChapterList", req.TitleID).Return(&cr, nil)
	s.ws.On("SetChapterCache", req.TitleID).Return(nil)

	cs := NewChapterService(s.cc, ccm, s.ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*chapters) > 0)
	assert.Equal(s.T(), "650", (*chapters)[0].Number)
	assert.Equal(s.T(), dc.Title, (*chapters)[0].Title)
	assert.Equal(s.T(), dc.TitleID, (*chapters)[0].TitleID)

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}
