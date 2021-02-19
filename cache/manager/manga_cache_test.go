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

type MangaCacheManagerTestSuite struct {
	suite.Suite
	mca cache.MangaCache
	mcl *mock.MangaClientMock
}

func TestMangaCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(MangaCacheManagerTestSuite))
}

func (s *MangaCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *MangaCacheManagerTestSuite) SetupTest() {
	s.mca = cache.NewMangaCache()
	s.mcl = &mock.MangaClientMock{}
}

func (s *MangaCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	s.mcl.On("GetMangaList").Return(nil, errors.New("some error"))

	mcm := NewMangaCacheManager(s.mcl, s.mca)
	err := mcm.SetCache()

	assert.Equal(s.T(), "some error", err.Error())
	s.mcl.AssertExpectations(s.T())
}

func (s *MangaCacheManagerTestSuite) TestSetCache_Succeed() {
	res := getFakeMangaList()
	s.mcl.On("GetMangaList").Return(&res, nil)

	mcm := NewMangaCacheManager(s.mcl, s.mca)
	err := mcm.SetCache()

	expCache, _ := json.Marshal(res)
	storedCache, _ := s.mca.Get()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(expCache), storedCache)
	s.mcl.AssertExpectations(s.T())

	_ = s.mca.Delete()
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	mcm := NewMangaCacheManager(s.mcl, s.mca)
	ml, err := mcm.GetCache()

	assert.Nil(s.T(), ml)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	mcm := NewMangaCacheManager(s.mcl, s.mca)

	_ = s.mca.Set("foo")
	defer func() {
		_ = s.mca.Delete()
	}()

	ml, err := mcm.GetCache()

	assert.Nil(s.T(), ml)
	assert.Equal(s.T(), "invalid manga cache", err.Error())
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsMangaList_WhenCacheIsStored() {
	mcm := NewMangaCacheManager(s.mcl, s.mca)

	cb, _ := json.Marshal(getFakeMangaList())
	_ = s.mca.Set(string(cb))
	defer func() {
		_ = s.mca.Delete()
	}()

	ml, err := mcm.GetCache()

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(ml.Mangas) > 0)
}

func getFakeMangaList() domain.MangaListResponse {
	return domain.MangaListResponse{
		Mangas: []domain.Manga{
			getFakePopularManga(),
			getFakeLatestManga(),
		},
	}
}

func getFakePopularManga() domain.Manga {
	return domain.Manga{
		ID:           "23",
		Title:        "One Piece",
		TitleID:      "one_piece",
		IconURL:      "http://foo.com/one_piece.jpg",
		LastChapter:  "939",
		ModifiedDate: "2019-04-12 13:05:59",
		Genre:        "adventure",
		Alias:        "n/a",
		Author:       "Eiichiro Oda",
		Status:       "ongoing",
		PublishYear:  "1997",
		Summary:      "Lorem ipsum...",
	}
}

func getFakeLatestManga() domain.Manga {
	return domain.Manga{
		ID:           "40",
		Title:        "Kagamigami",
		TitleID:      "kagamigami",
		IconURL:      "http://foo.com/kagamigami.jpg",
		LastChapter:  "30",
		ModifiedDate: "2019-04-12 11:26:49",
		Genre:        "fantasy",
		Alias:        "n/a",
		Author:       "Iwashiro Toshiaki",
		Status:       "ongoing",
		PublishYear:  "2015",
		Summary:      "Lorem ipsum...",
	}
}
