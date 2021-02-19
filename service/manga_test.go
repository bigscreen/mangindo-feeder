package service

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

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
)

type MangaServiceTestSuite struct {
	suite.Suite
	mca cache.MangaCache
	mc  *mock.MangaClientMock
	ws  *mock.WorkerServiceMock
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
	s.mc = &mock.MangaClientMock{}
	s.ws = &mock.WorkerServiceMock{}
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheMissesAndClientReturnsError() {
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	s.mc.On("GetMangaList").Return(nil, errors.New("some error"))

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewGenericError().Error(), err.Error())

	s.mc.AssertExpectations(s.T())
	s.ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheHitsAndMangaListIsEmpty() {
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	mr := domain.MangaListResponse{Mangas: []domain.Manga{}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer func() {
		_ = s.mca.Delete()
	}()

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewNotFoundError("manga").Error(), err.Error())

	s.mc.AssertNotCalled(s.T(), "GetMangaList")
	s.ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenCacheMissesAndMangaListIsEmpty() {
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	mr := domain.MangaListResponse{Mangas: []domain.Manga{}}
	s.mc.On("GetMangaList").Return(&mr, nil)
	s.ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewNotFoundError("manga").Error(), err.Error())

	s.mc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsPopularMangas_WhenCacheHitsAndMangaListContainsOnlyPopularMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	dm := getFakePopularManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer func() {
		_ = s.mca.Delete()
	}()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), lMangas)
	assert.True(s.T(), len(*pMangas) > 0)
	assertMappedManga(s.T(), dm, fpManga)

	s.mc.AssertNotCalled(s.T(), "GetMangaList")
	s.ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsPopularMangas_WhenCacheMissesAndMangaListContainsOnlyPopularMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dm := getFakePopularManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	s.mc.On("GetMangaList").Return(&mr, nil)
	s.ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), lMangas)
	assert.True(s.T(), len(*pMangas) > 0)
	assertMappedManga(s.T(), dm, fpManga)

	s.mc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsLatestMangas_WhenCacheHitsAndMangaListContainsOnlyLatestMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	dm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer func() {
		_ = s.mca.Delete()
	}()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), pMangas)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dm, flManga)

	s.mc.AssertNotCalled(s.T(), "GetMangaList")
	s.ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsLatestMangas_WhenCacheMissesAndMangaListContainsOnlyLatestMangas() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dm}}
	s.mc.On("GetMangaList").Return(&mr, nil)
	s.ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), pMangas)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dm, flManga)

	s.mc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsAllMangas_WhenCacheHits() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	dpm := getFakePopularManga()
	dlm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dpm, dlm}}
	cb, _ := json.Marshal(mr)
	_ = s.mca.Set(string(cb))
	defer func() {
		_ = s.mca.Delete()
	}()

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]
	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*pMangas) > 0)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dpm, fpManga)
	assertMappedManga(s.T(), dlm, flManga)

	s.mc.AssertNotCalled(s.T(), "GetMangaList")
	s.ws.AssertNotCalled(s.T(), "SetMangaCache")
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsAllMangas_WhenCacheMisses() {
	tags := os.Getenv("POPULAR_MANGA_TAGS")
	mcm := manager.NewMangaCacheManager(s.mc, s.mca)

	_ = os.Setenv("POPULAR_MANGA_TAGS", "one_piece")
	config.Load()
	defer func() {
		_ = os.Setenv("POPULAR_MANGA_TAGS", tags)
		config.Load()
	}()

	dpm := getFakePopularManga()
	dlm := getFakeLatestManga()
	mr := domain.MangaListResponse{Mangas: []domain.Manga{dpm, dlm}}
	s.mc.On("GetMangaList").Return(&mr, nil)
	s.ws.On("SetMangaCache").Return(nil)

	ms := NewMangaService(s.mc, mcm, s.ws)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]
	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*pMangas) > 0)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dpm, fpManga)
	assertMappedManga(s.T(), dlm, flManga)

	s.mc.AssertExpectations(s.T())
	s.ws.AssertExpectations(s.T())
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

func assertMappedManga(t *testing.T, dm domain.Manga, cm contract.Manga) {
	assert.Equal(t, dm.Title, cm.Title)
	assert.Equal(t, dm.TitleID, cm.TitleID)
	assert.Equal(t, dm.IconURL, cm.IconURL)
	assert.Equal(t, dm.LastChapter, cm.LastChapter)
	assert.Equal(t, dm.Genre, cm.Genre)
	assert.Equal(t, dm.Alias, cm.Alias)
	assert.Equal(t, dm.Author, cm.Author)
	assert.Equal(t, dm.Status, cm.Status)
	assert.Equal(t, dm.Summary, cm.Summary)
}
