package manager

import (
	"encoding/json"
	"errors"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ContentCacheManagerTestSuite struct {
	suite.Suite
}

func (s *ContentCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestContentCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCacheManagerTestSuite))
}

func (s *ContentCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	ccl := mock.MockContentClient{}
	cca := cache.NewContentCache()

	ccl.
		On("GetContentList", "bleach", float32(650.0)).
		Return(nil, errors.New("some error"))

	ccm := NewContentCacheManager(ccl, cca)
	err := ccm.SetCache("bleach", float32(650.0))

	assert.Equal(s.T(), "some error", err.Error())
	ccl.AssertExpectations(s.T())
}

func (s *ContentCacheManagerTestSuite) TestSetCache_WhenSucceed() {
	ccl := mock.MockContentClient{}
	cca := cache.NewContentCache()

	res := getFakeContentList()
	ccl.On("GetContentList", "bleach", float32(650.0)).Return(&res, nil)

	ccm := NewContentCacheManager(ccl, cca)
	err := ccm.SetCache("bleach", float32(650.0))

	ec, _ := json.Marshal(res)
	sc, _ := cca.Get("bleach", "650")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(ec), sc)
	ccl.AssertExpectations(s.T())

	_ = cca.Delete("bleach", "650")
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	cca := cache.NewContentCache()
	ccm := NewContentCacheManager(mock.MockContentClient{}, cca)

	cl, err := ccm.GetCache("bleach", float32(650.0))

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	cca := cache.NewContentCache()
	ccm := NewContentCacheManager(mock.MockContentClient{}, cca)

	_ = cca.Set("bleach", "650", "foo")
	defer cca.Delete("bleach", "650")

	cl, err := ccm.GetCache("bleach", float32(650.0))

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "invalid content cache", err.Error())
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsContentList_WhenCacheIsStored() {
	cca := cache.NewContentCache()
	ccm := NewContentCacheManager(mock.MockContentClient{}, cca)

	cb, _ := json.Marshal(getFakeContentList())
	_ = cca.Set("bleach", "650", string(cb))
	defer cca.Delete("bleach", "650")

	cl, err := ccm.GetCache("bleach", float32(650.0))

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(cl.Contents) > 0)
}

func getFakeContentList() domain.ContentListResponse {
	return domain.ContentListResponse{
		Contents: []domain.Content{
			getFakeContent(1),
			getFakeContent(2),
		},
	}
}

func getFakeContent(page int) domain.Content {
	return domain.Content{
		ImageURL: "http://foo.com/foo.jpg",
		Page:     page,
	}
}
