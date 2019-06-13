package service

import (
	"encoding/json"
	"errors"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/cache/manager"
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

type MangaServiceTestSuite struct {
	suite.Suite
	mca cache.MangaCache
}

func TestMangaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MangaServiceTestSuite))
}

func (s *MangaServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *MangaServiceTestSuite) SetupTest() {
	s.mca = cache.NewMangaCache()
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	mc.On("GetMangaList").Return(nil, errors.New("some error"))

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	mc.AssertExpectations(s.T())
	ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheHitsAndMangaListIsEmpty() {
	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	mr := domain.MangaListResponse{Mangas: []domain.Manga{}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer s.mca.Delete()

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewNotFoundError("manga").Error(), err.Error())

	mc.AssertNotCalled(s.T(), "GetMangaList")
	ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheMissesAndMangaListIsEmpty() {
	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	mr := domain.MangaListResponse{Mangas: []domain.Manga{}}
	mc.On("GetMangaList").Return(&mr, nil)
	ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewNotFoundError("manga").Error(), err.Error())

	mc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsPopularMangas_WhenCacheHitsAndMangaListContainsOnlyPopularMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	dm := getFakePopularManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer s.mca.Delete()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), lMangas)
	assert.True(s.T(), len(*pMangas) > 0)
	assertMappedManga(s.T(), dm, fpManga)

	mc.AssertNotCalled(s.T(), "GetMangaList")
	ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsPopularMangas_WhenCacheMissesAndMangaListContainsOnlyPopularMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dm := getFakePopularManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	mc.On("GetMangaList").Return(&mr, nil)
	ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), lMangas)
	assert.True(s.T(), len(*pMangas) > 0)
	assertMappedManga(s.T(), dm, fpManga)

	mc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsLatestMangas_WhenCacheHitsAndMangaListContainsOnlyLatestMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	dm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer s.mca.Delete()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), pMangas)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dm, flManga)

	mc.AssertNotCalled(s.T(), "GetMangaList")
	ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsLatestMangas_WhenCacheMissesAndMangaListContainsOnlyLatestMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	mc.On("GetMangaList").Return(&mr, nil)
	ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), pMangas)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dm, flManga)

	mc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsAllMangas_WhenCacheHits() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	dpm := getFakePopularManga()
	dlm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dpm, dlm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer s.mca.Delete()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]
	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*pMangas) > 0)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dpm, fpManga)
	assertMappedManga(s.T(), dlm, flManga)

	mc.AssertNotCalled(s.T(), "GetMangaList")
	ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsAllMangas_WhenCacheMisses() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")

	mc := mock.MangaClientMock{}
	ws := mock.WorkerServiceMock{}
	mcm := manager.NewMangaCacheManager(mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dpm := getFakePopularManga()
	dlm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dpm, dlm}}
	mc.On("GetMangaList").Return(&mr, nil)
	ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(mc, mcm, ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]
	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*pMangas) > 0)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dpm, fpManga)
	assertMappedManga(s.T(), dlm, flManga)

	mc.AssertExpectations(s.T())
	ws.AssertExpectations(s.T())
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

func assertMappedManga(t *testing.T, dm domain.Manga, cm contract.Manga) {
	assert.Equal(t, dm.Title, cm.Title)
	assert.Equal(t, dm.TitleId, cm.TitleId)
	assert.Equal(t, dm.IconURL, cm.IconURL)
	assert.Equal(t, dm.LastChapter, cm.LastChapter)
	assert.Equal(t, dm.Genre, cm.Genre)
	assert.Equal(t, dm.Alias, cm.Alias)
	assert.Equal(t, dm.Author, cm.Author)
	assert.Equal(t, dm.Status, cm.Status)
	assert.Equal(t, dm.Summary, cm.Summary)
}
