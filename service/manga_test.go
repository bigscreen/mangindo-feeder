package service

import (
	"errors"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MangaServiceTestSuite struct {
	suite.Suite
}

func (s *MangaServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestMangaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MangaServiceTestSuite))
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenClientReturnsError() {
	mc := client.MockMangaClient{}

	mc.On("GetMangaList").Return(nil, errors.New("some error"))

	ms := NewMangaService(mc)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), "some error", err.Error())
	mc.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsError_WhenMangaListIsEmpty() {
	mc := client.MockMangaClient{}
	res := &domain.MangaListResponse{Mangas: []domain.Manga{},}

	mc.On("GetMangaList").Return(res, nil)

	ms := NewMangaService(mc)
	pMangas, lMangas, err := ms.GetMangas()

	assert.Nil(s.T(), pMangas)
	assert.Nil(s.T(), lMangas)
	assert.Equal(s.T(), mErr.NewNotFoundError("manga").Error(), err.Error())
	mc.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsOnlyPopularMangas_WhenMangaListContainsOnlyPopularMangas() {
	mc := client.MockMangaClient{}
	dm := getFakePopularManga()
	res := &domain.MangaListResponse{Mangas: []domain.Manga{dm}}

	mc.On("GetMangaList").Return(res, nil)

	ms := NewMangaService(mc)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), lMangas)
	assert.True(s.T(), len(*pMangas) > 0)
	assertMappedManga(s.T(), dm, fpManga)
	mc.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsOnlyLatestMangas_WhenMangaListContainsOnlyLatestMangas() {
	mc := client.MockMangaClient{}
	dm := getFakeLatestManga()
	res := &domain.MangaListResponse{Mangas: []domain.Manga{dm}}

	mc.On("GetMangaList").Return(res, nil)

	ms := NewMangaService(mc)
	pMangas, lMangas, err := ms.GetMangas()

	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.Nil(s.T(), pMangas)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dm, flManga)
	mc.AssertExpectations(s.T())
}

func (s *MangaServiceTestSuite) TestGetMangas_ReturnsAllMangas() {
	mc := client.MockMangaClient{}
	dpm := getFakePopularManga()
	dlm := getFakeLatestManga()
	res := &domain.MangaListResponse{Mangas: []domain.Manga{dpm, dlm}}

	mc.On("GetMangaList").Return(res, nil)

	ms := NewMangaService(mc)
	pMangas, lMangas, err := ms.GetMangas()

	fpManga := (*pMangas)[0]
	flManga := (*lMangas)[0]

	assert.Nil(s.T(), err)
	assert.True(s.T(), len(*pMangas) > 0)
	assert.True(s.T(), len(*lMangas) > 0)
	assertMappedManga(s.T(), dpm, fpManga)
	assertMappedManga(s.T(), dlm, flManga)
	mc.AssertExpectations(s.T())
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
