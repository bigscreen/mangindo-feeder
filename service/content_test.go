package service

import (
	"errors"
	"fmt"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
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
	cc := mock.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(nil, errors.New("some error"))

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenContentListIsEmpty() {
	cc := mock.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	res := &domain.ContentListResponse{
		Contents: []domain.Content{},
	}

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenContentListContainsOnlyAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	res := &domain.ContentListResponse{
		Contents: []domain.Content{getFakeAdsContent(1, "ads")},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenContentListContainsOnlyNonAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	ct := getFakeContent(1)
	res := &domain.ContentListResponse{
		Contents: []domain.Content{ct},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct.ImageURL), (*cl)[0].ImageURL)
	cc.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenContentListContainsAdsAndNonAdsContents() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.MockContentClient{}
	req := contract.NewContentRequest("bleach", "650")
	ct1 := getFakeAdsContent(1, "ads")
	ct2 := getFakeContent(2)
	res := &domain.ContentListResponse{
		Contents: []domain.Content{ct1, ct2},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(res, nil)

	cs := NewContentService(cc)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct2.ImageURL), (*cl)[0].ImageURL)
	cc.AssertExpectations(s.T())
}

func getFakeContent(page int) domain.Content {
	return domain.Content{
		ImageURL: "http://foo.com/foo.jpg",
		Page:     page,
	}
}

func getFakeAdsContent(page int, adsTag string) domain.Content {
	return domain.Content{
		ImageURL: fmt.Sprintf("http://foo.com/%s.jpg", adsTag),
		Page:     page,
	}
}
