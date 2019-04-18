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

type ContentServiceTestSuite struct {
	suite.Suite
}

func (s *ContentServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestContentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContentServiceTestSuite))
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenClientReturnsError() {
	cc := client.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(nil, errors.New("some error"))

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenContentListIsEmpty() {
	cc := client.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	res := &domain.ContentListResponse{Contents: []domain.Content{},}

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenContentListIsEmpty() {
	cc := client.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	ct := domain.Content{
		ImageURL: "http://foo.com/foo.jpg",
		Page:     1,
	}
	res := &domain.ContentListResponse{Contents: []domain.Content{ct},}

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct.ImageURL), (*cl)[0].ImageURL)
	cc.AssertExpectations(s.T())
}
