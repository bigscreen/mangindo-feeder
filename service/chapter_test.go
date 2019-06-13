package service

import (
	"encoding/json"
	"errors"
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
	"testing"
)

type ChapterServiceTestSuite struct {
	suite.Suite
	cca cache.ChapterCache
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
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	cc := mock.ChapterClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewChapterCacheManager(cc, s.cca)

	req := contract.NewChapterRequest("bleach")

	cc.On("GetChapterList", req.TitleId).Return(nil, errors.New("some error"))

	cs := NewChapterService(cc, ccm, ws)
	cl, err := cs.GetChapters(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	cc.AssertExpectations(s.T())
	ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleId)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheHitsAndChapterListIsEmpty() {
	cc := mock.ChapterClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewChapterCacheManager(cc, s.cca)

	req := contract.NewChapterRequest("bleach")

	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{}}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, string(cb))
	defer s.cca.Delete(req.TitleId)

	cs := NewChapterService(cc, ccm, ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), chapters)
	assert.Equal(s.T(), mErr.NewNotFoundError("chapter").Error(), err.Error())

	ws.AssertNotCalled(s.T(), "GetChapterList", req.TitleId)
	ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleId)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenCacheMissesAndChapterListIsEmpty() {
	cc := mock.ChapterClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewChapterCacheManager(cc, s.cca)

	req := contract.NewChapterRequest("bleach")

	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{}}

	cc.On("GetChapterList", req.TitleId).Return(&cr, nil)
	ws.On("SetChapterCache", req.TitleId).Return(nil)

	cs := NewChapterService(cc, ccm, ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), chapters)
	assert.Equal(s.T(), mErr.NewNotFoundError("chapter").Error(), err.Error())

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsSuccess_WhenCacheHitsAndChapterListIsNotEmpty() {
	cc := mock.ChapterClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewChapterCacheManager(cc, s.cca)

	req := contract.NewChapterRequest("bleach")

	dc := domain.Chapter{
		Number:       650.0,
		Title:        "Bleach",
		TitleId:      "bleach",
		ModifiedDate: "2016-08-18 18:59:58",
	}
	cr := &domain.ChapterListResponse{Chapters: []domain.Chapter{dc}}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, string(cb))
	defer s.cca.Delete(req.TitleId)

	cs := NewChapterService(cc, ccm, ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*chapters) > 0)
	assert.Equal(s.T(), "650", (*chapters)[0].Number)
	assert.Equal(s.T(), dc.Title, (*chapters)[0].Title)
	assert.Equal(s.T(), dc.TitleId, (*chapters)[0].TitleId)

	ws.AssertNotCalled(s.T(), "GetChapterList", req.TitleId)
	ws.AssertNotCalled(s.T(), "SetChapterCache", req.TitleId)
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsSuccess_WhenCacheMissesAndChapterListIsNotEmpty() {
	cc := mock.ChapterClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewChapterCacheManager(cc, s.cca)

	req := contract.NewChapterRequest("bleach")
	dc := domain.Chapter{
		Number:       650.0,
		Title:        "Bleach",
		TitleId:      "bleach",
		ModifiedDate: "2016-08-18 18:59:58",
	}
	cr := domain.ChapterListResponse{Chapters: []domain.Chapter{dc}}

	cc.On("GetChapterList", req.TitleId).Return(&cr, nil)
	ws.On("SetChapterCache", req.TitleId).Return(nil)

	cs := NewChapterService(cc, ccm, ws)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*chapters) > 0)
	assert.Equal(s.T(), "650", (*chapters)[0].Number)
	assert.Equal(s.T(), dc.Title, (*chapters)[0].Title)
	assert.Equal(s.T(), dc.TitleId, (*chapters)[0].TitleId)

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}
