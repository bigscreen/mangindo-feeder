package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/common"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContentServiceTestSuite struct {
	suite.Suite
	cca cache.ContentCache
	cc  *mock.ContentClientMock
	ws  *mock.WorkerServiceMock
}

func TestContentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ContentServiceTestSuite))
}

func (s *ContentServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ContentServiceTestSuite) SetupTest() {
	s.cca = cache.NewContentCache()
	s.cc = &mock.ContentClientMock{}
	s.ws = &mock.WorkerServiceMock{}
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	s.cc.On("GetContentList", req.TitleID, req.Chapter).Return(nil, errors.New("some error"))

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	s.cc.AssertExpectations(s.T())
	s.ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleID, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheHitsAndContentListIsEmpty() {
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, sch, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID, sch)
	}()

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	s.cc.AssertNotCalled(s.T(), "GetContentList", req.TitleID, req.Chapter)
	s.ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleID, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndContentListIsEmpty() {
	ccm := manager.NewContentCacheManager(s.cc, s.cca)
	req := contract.NewContentRequest("bleach", "650")
	cr := domain.ContentListResponse{Contents: []domain.Content{}}

	s.cc.On("GetContentList", req.TitleID, req.Chapter).Return(&cr, nil)
	s.ws.On("SetContentCache", req.TitleID, req.Chapter).Return(nil)

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheHitsAndContentListContainsOnlyAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{getFakeAdsContent(1, "ads")},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, sch, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID, sch)
	}()

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	s.cc.AssertNotCalled(s.T(), "GetContentList", req.TitleID, req.Chapter)
	s.ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleID, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndContentListContainsOnlyAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	cr := domain.ContentListResponse{
		Contents: []domain.Content{getFakeAdsContent(1, "ads")},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	s.cc.On("GetContentList", req.TitleID, req.Chapter).Return(&cr, nil)
	s.ws.On("SetContentCache", req.TitleID, req.Chapter).Return(nil)

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheHitsAndContentListContainsOnlyNonAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)
	ct := getFakeContent(1)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, sch, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID, sch)
	}()

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedURL(ct.ImageURL), (*cl)[0].ImageURL)

	s.cc.AssertNotCalled(s.T(), "GetContentList", req.TitleID, req.Chapter)
	s.ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleID, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheMissesAndContentListContainsOnlyNonAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	ct := getFakeContent(1)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	s.cc.On("GetContentList", req.TitleID, req.Chapter).Return(&cr, nil)
	s.ws.On("SetContentCache", req.TitleID, req.Chapter).Return(nil)

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedURL(ct.ImageURL), (*cl)[0].ImageURL)

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheHitsAndContentListContainsAdsAndNonAdsContents() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)
	ct1 := getFakeAdsContent(1, "ads")
	ct2 := getFakeContent(2)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct1, ct2},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleID, sch, string(cb))
	defer func() {
		_ = s.cca.Delete(req.TitleID, sch)
	}()

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedURL(ct2.ImageURL), (*cl)[0].ImageURL)

	s.cc.AssertNotCalled(s.T(), "GetContentList", req.TitleID, req.Chapter)
	s.ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleID, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenContentListContainsAdsAndNonAdsContents() {
	tags := os.Getenv("ADS_CONTENT_TAGS")
	ccm := manager.NewContentCacheManager(s.cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	ct1 := getFakeAdsContent(1, "ads")
	ct2 := getFakeContent(2)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct1, ct2},
	}

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	s.cc.On("GetContentList", req.TitleID, req.Chapter).Return(&cr, nil)
	s.ws.On("SetContentCache", req.TitleID, req.Chapter).Return(nil)

	cs := NewContentService(s.cc, ccm, s.ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedURL(ct2.ImageURL), (*cl)[0].ImageURL)

	s.cc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
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
