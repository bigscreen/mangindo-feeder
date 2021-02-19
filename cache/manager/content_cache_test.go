package manager

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContentCacheManagerTestSuite struct {
	suite.Suite
	cca cache.ContentCache
	ccl *mock.ContentClientMock
}

func TestContentCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCacheManagerTestSuite))
}

func (s *ContentCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ContentCacheManagerTestSuite) SetupTest() {
	s.cca = cache.NewContentCache()
	s.ccl = &mock.ContentClientMock{}
}

func (s *ContentCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	s.ccl.On("GetContentList", "bleach", float32(650.0)).
		Return(nil, errors.New("some error"))

	ccm := NewContentCacheManager(s.ccl, s.cca)
	err := ccm.SetCache("bleach", float32(650.0))

	assert.Equal(s.T(), "some error", err.Error())
	s.ccl.AssertExpectations(s.T())
}

func (s *ContentCacheManagerTestSuite) TestSetCache_WhenSucceed() {
	res := getFakeContentList()
	s.ccl.On("GetContentList", "bleach", float32(650.0)).Return(&res, nil)

	ccm := NewContentCacheManager(s.ccl, s.cca)
	err := ccm.SetCache("bleach", float32(650.0))

	ec, _ := json.Marshal(res)
	sc, _ := s.cca.Get("bleach", "650")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(ec), sc)
	s.ccl.AssertExpectations(s.T())

	_ = s.cca.Delete("bleach", "650")
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	ccm := NewContentCacheManager(s.ccl, s.cca)

	cl, err := ccm.GetCache("bleach", float32(650.0))

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	ccm := NewContentCacheManager(s.ccl, s.cca)

	_ = s.cca.Set("bleach", "650", "foo")
	defer func() {
		_ = s.cca.Delete("bleach", "650")
	}()

	cl, err := ccm.GetCache("bleach", float32(650.0))

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "invalid content cache", err.Error())
}

func (s *ContentCacheManagerTestSuite) TestGetCache_ReturnsContentList_WhenCacheIsStored() {
	ccm := NewContentCacheManager(s.ccl, s.cca)

	cb, _ := json.Marshal(getFakeContentList())
	_ = s.cca.Set("bleach", "650", string(cb))
	defer func() {
		_ = s.cca.Delete("bleach", "650")
	}()

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
