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

type MangaCacheManagerTestSuite struct {
	suite.Suite
}

func (s *MangaCacheManagerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestMangaCacheManagerTestSuite(t *testing.T) {
	suite.Run(t, new(MangaCacheManagerTestSuite))
}

func (s *MangaCacheManagerTestSuite) TestSetCache_ReturnsError_WhenClientReturnsError() {
	mcl := mock.MockMangaClient{}
	mca := cache.NewMangaCache()

	mcl.On("GetMangaList").Return(nil, errors.New("some error"))

	mcm := NewMangaCacheManager(mcl, mca)
	err := mcm.SetCache()

	assert.Equal(s.T(), "some error", err.Error())
	mcl.AssertExpectations(s.T())
}

func (s *MangaCacheManagerTestSuite) TestSetCache_Succeed() {
	mcl := mock.MockMangaClient{}
	mca := cache.NewMangaCache()

	res := getFakeMangaList()
	mcl.On("GetMangaList").Return(&res, nil)

	mcm := NewMangaCacheManager(mcl, mca)
	err := mcm.SetCache()

	expCache, _ := json.Marshal(res)
	storedCache, _ := mca.Get()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(expCache), storedCache)
	mcl.AssertExpectations(s.T())

	_ = mca.Delete()
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsMissing() {
	mca := cache.NewMangaCache()
	mcm := NewMangaCacheManager(mock.MockMangaClient{}, mca)

	ml, err := mcm.GetCache()

	assert.Nil(s.T(), ml)
	assert.Equal(s.T(), "redis: nil", err.Error())
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsError_WhenCacheIsInvalid() {
	mca := cache.NewMangaCache()
	mcm := NewMangaCacheManager(mock.MockMangaClient{}, mca)

	_ = mca.Set("foo")
	defer mca.Delete()

	ml, err := mcm.GetCache()

	assert.Nil(s.T(), ml)
	assert.Equal(s.T(), "invalid manga cache", err.Error())
}

func (s *MangaCacheManagerTestSuite) TestGetCache_ReturnsMangaList_WhenCacheIsStored() {
	mca := cache.NewMangaCache()
	mcm := NewMangaCacheManager(mock.MockMangaClient{}, mca)

	cb, _ := json.Marshal(getFakeMangaList())
	_ = mca.Set(string(cb))
	defer mca.Delete()

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
		Id:           "23",
		Title:        "One Piece",
		TitleId:      "one_piece",
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
		Id:           "40",
		Title:        "Kagamigami",
		TitleId:      "kagamigami",
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
