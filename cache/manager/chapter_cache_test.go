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

type ChapterCacheManagerTestSuite struct {
	suite.Suite
}

func (s *ChapterCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestChapterCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterCacheManagerTestSuite))
}

func (s *ChapterCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	ccl := mock.ChapterClientMock{}
	cca := cache.NewChapterCache()

	ccl.On("GetChapterList", "bleach").Return(nil, errors.New("some error"))

	ccm := NewChapterCacheManager(ccl, cca)
	err := ccm.SetCache("bleach")

	assert.Equal(s.T(), "some error", err.Error())
	ccl.AssertExpectations(s.T())
}

func (s *ChapterCacheManagerTestSuite) TestSetCache_WhenSucceed() {
	ccl := mock.ChapterClientMock{}
	cca := cache.NewChapterCache()

	res := getFakeChapterList()
	ccl.On("GetChapterList", "bleach").Return(&res, nil)

	ccm := NewChapterCacheManager(ccl, cca)
	err := ccm.SetCache("bleach")

	ec, _ := json.Marshal(res)
	sc, _ := cca.Get("bleach")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(ec), sc)
	ccl.AssertExpectations(s.T())

	_ = cca.Delete("bleach")
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	cca := cache.NewChapterCache()
	ccm := NewChapterCacheManager(mock.ChapterClientMock{}, cca)

	cl, err := ccm.GetCache("bleach")

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	cca := cache.NewChapterCache()
	ccm := NewChapterCacheManager(mock.ChapterClientMock{}, cca)

	_ = cca.Set("bleach", "foo")
	defer cca.Delete("bleach")

	cl, err := ccm.GetCache("bleach")

	assert.Nil(s.T(), cl)
	assert.Equal(s.T(), "invalid chapter cache", err.Error())
}

func (s *ChapterCacheManagerTestSuite) TestGetCache_ReturnsChapterList_WhenCacheIsStored() {
	cca := cache.NewChapterCache()
	ccm := NewChapterCacheManager(mock.ChapterClientMock{}, cca)

	cb, _ := json.Marshal(getFakeChapterList())
	_ = cca.Set("bleach", string(cb))
	defer cca.Delete("bleach")

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
				TitleId:      "bleach",
				ModifiedDate: "2016-08-18 18:59:58",
			},
		},
	}
}
