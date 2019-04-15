package service

import (
	"errors"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ChapterServiceTestSuite struct {
	suite.Suite
}

func (s *ChapterServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestChapterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterServiceTestSuite))
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenClientReturnsError() {
	cc := client.MockChapterClient{}
	req := contract.NewChapterRequest("bleach")

	cc.On("GetChapterList", req.TitleId).Return(nil, errors.New("some error"))

	cs := NewChapterService(cc)
	cl, err := cs.GetChapters(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsError_WhenChapterListIsEmpty() {
	cc := client.MockChapterClient{}
	req := contract.NewChapterRequest("bleach")
	res := &domain.ChapterListResponse{Chapters: []domain.Chapter{},}

	cc.On("GetChapterList", req.TitleId).Return(res, nil)

	cs := NewChapterService(cc)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), chapters)
	assert.Equal(s.T(), mErr.NewNotFoundError("chapter").Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ChapterServiceTestSuite) TestGetChapters_ReturnsSuccess_WhenChapterListIsNotEmpty() {
	cc := client.MockChapterClient{}
	req := contract.NewChapterRequest("bleach")
	dc := domain.Chapter{
		Number:       650.0,
		Title:        "Bleach",
		TitleId:      "bleach",
		ModifiedDate: "2016-08-18 18:59:58",
	}
	res := &domain.ChapterListResponse{Chapters: []domain.Chapter{dc},}

	cc.On("GetChapterList", req.TitleId).Return(res, nil)

	cs := NewChapterService(cc)
	chapters, err := cs.GetChapters(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*chapters) > 0)
	assert.Equal(s.T(), "650", (*chapters)[0].Number)
	assert.Equal(s.T(), dc.Title, (*chapters)[0].Title)
	assert.Equal(s.T(), dc.TitleId, (*chapters)[0].TitleId)
	cc.AssertExpectations(s.T())
}
