package service

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"os"
	"testing"
)

type ContentServiceTestSuite struct {
	suite.Suite
	cca cache.ContentCache
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
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(nil, errors.New("some error"))

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	cc.AssertExpectations(s.T())
	ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleId, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheHitsAndContentListIsEmpty() {
	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)

	cr := domain.ContentListResponse{
		Contents: []domain.Content{},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, sch, string(cb))
	defer s.cca.Delete(req.TitleId, sch)

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	cc.AssertNotCalled(s.T(), "GetContentList", req.TitleId, req.Chapter)
	ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleId, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndContentListIsEmpty() {
	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")

	cr := domain.ContentListResponse{Contents: []domain.Content{}}

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(&cr, nil)
	ws.On("SetContentCache", req.TitleId, req.Chapter).Return(nil)

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheHitsAndContentListContainsOnlyAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)

	cr := domain.ContentListResponse{
		Contents: []domain.Content{getFakeAdsContent(1, "ads")},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, sch, string(cb))
	defer s.cca.Delete(req.TitleId, sch)

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	cc.AssertNotCalled(s.T(), "GetContentList", req.TitleId, req.Chapter)
	ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleId, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsError_WhenCacheMissesAndContentListContainsOnlyAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

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

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(&cr, nil)
	ws.On("SetContentCache", req.TitleId, req.Chapter).Return(nil)

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), mErr.NewNotFoundError("content").Error(), err.Error())

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheHitsAndContentListContainsOnlyNonAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)

	ct := getFakeContent(1)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, sch, string(cb))
	defer s.cca.Delete(req.TitleId, sch)

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct.ImageURL), (*cl)[0].ImageURL)

	cc.AssertNotCalled(s.T(), "GetContentList", req.TitleId, req.Chapter)
	ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleId, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheMissesAndContentListContainsOnlyNonAdsContent() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

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

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(&cr, nil)
	ws.On("SetContentCache", req.TitleId, req.Chapter).Return(nil)

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct.ImageURL), (*cl)[0].ImageURL)

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenCacheHitsAndContentListContainsAdsAndNonAdsContents() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

	req := contract.NewContentRequest("bleach", "650")
	sch := common.GetFormattedChapterNumber(req.Chapter)

	ct1 := getFakeAdsContent(1, "ads")
	ct2 := getFakeContent(2)
	cr := domain.ContentListResponse{
		Contents: []domain.Content{ct1, ct2},
	}
	cb, _ := json.Marshal(cr)
	_ = s.cca.Set(req.TitleId, sch, string(cb))
	defer s.cca.Delete(req.TitleId, sch)

	_ = os.Setenv("ADS_CONTENT_TAGS", "ads")
	config.Load()
	defer func() {
		_ = os.Setenv("ADS_CONTENT_TAGS", tags)
		config.Load()
	}()

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct2.ImageURL), (*cl)[0].ImageURL)

	cc.AssertNotCalled(s.T(), "GetContentList", req.TitleId, req.Chapter)
	ws.AssertNotCalled(s.T(), "SetContentCache", req.TitleId, req.Chapter)
}

func (s *ContentServiceTestSuite) TestGetContents_ReturnsSuccess_WhenContentListContainsAdsAndNonAdsContents() {
	tags := os.Getenv("ADS_CONTENT_TAGS")

	cc := mock.ContentClientMock{}
	ws := mock.WorkerServiceMock{}
	ccm := manager.NewContentCacheManager(cc, s.cca)

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

	cc.On("GetContentList", req.TitleId, req.Chapter).Return(&cr, nil)
	ws.On("SetContentCache", req.TitleId, req.Chapter).Return(nil)

	cs := NewContentService(cc, ccm, ws)
	cl, err := cs.GetContents(req)

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*cl) > 0)
	assert.Equal(s.T(), getEncodedUrl(ct2.ImageURL), (*cl)[0].ImageURL)

	cc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
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
