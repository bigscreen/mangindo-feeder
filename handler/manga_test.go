package handler

import (
	"encoding/json"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	mMock "github.com/bigscreen/mangindo-feeder/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MangaHandlerTestSuite struct {
	suite.Suite
}

func (s *MangaHandlerTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func TestMangaHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MangaHandlerTestSuite))
}

func (s *MangaHandlerTestSuite) TestGetMangas_ReturnsError_WhenUnknownErrorHappens() {
	req, _ := http.NewRequest("GET", constants.GetMangasApiPath, nil)

	rr := httptest.NewRecorder()
	err := mErr.NewGenericError()

	ms := mMock.MockMangaService{}
	ms.On("GetMangas").Return(nil, nil, err)

	h := GetMangas(ms)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusInternalServerError, rr.Code)
	assert.Equal(s.T(), string(res), body)
	ms.AssertExpectations(s.T())
}

func (s *MangaHandlerTestSuite) TestGetMangas_ReturnsError_WhenMangasDoNotExist() {
	req, _ := http.NewRequest("GET", constants.GetMangasApiPath, nil)

	rr := httptest.NewRecorder()
	err := mErr.NewNotFoundError("manga")

	ms := mMock.MockMangaService{}
	ms.On("GetMangas").Return(nil, nil, err)

	h := GetMangas(ms)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(getErrorResponse(err))
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusNotFound, rr.Code)
	assert.Equal(s.T(), string(res), body)
	ms.AssertExpectations(s.T())
}

func (s *MangaHandlerTestSuite) TestGetMangas_ReturnsSuccess_WhenOnlyPopularMangasExist() {
	req, _ := http.NewRequest("GET", constants.GetMangasApiPath, nil)

	rr := httptest.NewRecorder()
	pms := []contract.Manga{getFakePopularManga()}
	mr := contract.MangaResponse{
		Success:       true,
		PopularMangas: pms,
		LatestMangas:  []contract.Manga{},
	}

	ms := mMock.MockMangaService{}
	ms.On("GetMangas").Return(&pms, nil, nil)

	h := GetMangas(ms)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(mr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	ms.AssertExpectations(s.T())
}

func (s *MangaHandlerTestSuite) TestGetMangas_ReturnsSuccess_WhenOnlyLatestMangasExist() {
	req, _ := http.NewRequest("GET", constants.GetMangasApiPath, nil)

	rr := httptest.NewRecorder()
	lms := []contract.Manga{getFakeLatestManga()}
	mr := contract.MangaResponse{
		Success:       true,
		PopularMangas: []contract.Manga{},
		LatestMangas:  lms,
	}

	ms := mMock.MockMangaService{}
	ms.On("GetMangas").Return(nil, &lms, nil)

	h := GetMangas(ms)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(mr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	ms.AssertExpectations(s.T())
}

func (s *MangaHandlerTestSuite) TestGetMangas_ReturnsSuccess_WhenPopularAndLatestMangasExist() {
	req, _ := http.NewRequest("GET", constants.GetMangasApiPath, nil)

	rr := httptest.NewRecorder()
	pms := []contract.Manga{getFakePopularManga()}
	lms := []contract.Manga{getFakeLatestManga()}
	mr := contract.MangaResponse{
		Success:       true,
		PopularMangas: pms,
		LatestMangas:  lms,
	}

	ms := mMock.MockMangaService{}
	ms.On("GetMangas").Return(&pms, &lms, nil)

	h := GetMangas(ms)
	h.ServeHTTP(rr, req)

	res, _ := json.Marshal(mr)
	body := strings.TrimSuffix(string(rr.Body.Bytes()), "\n")

	assert.Equal(s.T(), http.StatusOK, rr.Code)
	assert.Equal(s.T(), string(res), body)
	ms.AssertExpectations(s.T())
}

func getFakePopularManga() contract.Manga {
	return contract.Manga{
		Title:       "One Piece",
		TitleId:     "one_piece",
		IconURL:     "http://foo.com/one_piece.jpg",
		LastChapter: "939",
		Genre:       "adventure",
		Alias:       "n/a",
		Author:      "Eiichiro Oda",
		Status:      "ongoing",
		PublishYear: "1997",
		Summary:     "Lorem ipsum...",
	}
}

func getFakeLatestManga() contract.Manga {
	return contract.Manga{
		Title:       "Kagamigami",
		TitleId:     "kagamigami",
		IconURL:     "http://foo.com/kagamigami.jpg",
		LastChapter: "30",
		Genre:       "fantasy",
		Alias:       "n/a",
		Author:      "Iwashiro Toshiaki",
		Status:      "ongoing",
		PublishYear: "2015",
		Summary:     "Lorem ipsum...",
	}
}
