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

type ChapterCacheManagerTestSuite struct {
	suite.Suite
	cca cache.ChapterCache
	ccl *mock.ChapterClientMock
}

func TestChapterCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterCacheManagerTestSuite))
}

func (s *ChapterCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *ChapterCacheManagerTestSuite) SetupTest() {
	s.cca = cache.NewChapterCache()
	s.ccl = &mock.ChapterClientMock{}
}

func (s *ChapterCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	s.ccl.On("GetChapterList", "bleach").Return(nil, errors.New("some error"))

	ccm := NewChapterCacheManager(s.ccl, s.cca)
	err := ccm.SetCache("bleach")

	assert.Equal(s.T(), "some error", err.Error())
	s.ccl.AssertExpectations(s.T())
}

func (s *ChapterCacheManagerTestSuite) TestSetCache_WhenSucceed() {
	res := getFakeChapterList()
	s.ccl.On("GetChapterList", "bleach").Return(&res, nil)

	ccm := NewChapterCacheManager(s.ccl, s.cca)
	err := ccm.SetCache("bleach")

	ec, _ := json.Marshal(res)
	sc, _ := s.cca.Get("bleach")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(ec), sc)
	s.ccl.AssertExpectations(s.T())

	_ = s.cca.Delete("bleach")
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	ccm := NewChapterCacheManager(s.ccl, s.cca)

	cl, err := ccm.GetCache("bleach")

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	ccm := NewChapterCacheManager(s.ccl, s.cca)

	_ = s.cca.Set("bleach", "foo")
	defer func() {
		_ = s.cca.Delete("bleach")
	}()

	cl, err := ccm.GetCache("bleach")

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "invalid chapter cache", err.Error())
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsChapterList_WhenCacheIsStored() {
	ccm := NewChapterCacheManager(s.ccl, s.cca)

	cb, _ := json.Marshal(getFakeChapterList())
	_ = s.cca.Set("bleach", string(cb))
	defer func() {
		_ = s.cca.Delete("bleach")
	}()

	cl, err := ccm.GetCache("bleach")

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(cl.Chapters) > 0)
}

func getFakeChapterList() domain.ChapterListResponse {
	return domain.ChapterListResponse{
		Chapters: []domain.Chapter{
			{
				Number:       650.0,
				Title:        "Bleach",
				TitleID:      "bleach",
				ModifiedDate: "2016-08-18 18:59:58",
			},
		},
	}
}
