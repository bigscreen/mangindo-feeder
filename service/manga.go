package service

import (
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
)

type MangaService interface {
	GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error)
}

type mangaService struct {
	mClient client.MangaClient
}

func getMappedManga(dm domain.Manga) contract.Manga {
	return contract.Manga{
		Title:       dm.Title,
		TitleId:     dm.TitleId,
		IconURL:     dm.IconURL,
		LastChapter: dm.LastChapter,
		Genre:       dm.Genre,
		Alias:       dm.Alias,
		Author:      dm.Author,
		Status:      dm.Status,
		PublishYear: dm.PublishYear,
		Summary:     dm.Summary,
	}
}

func isPopularManga(titleId string) bool {
	popularTags := [5]string{"one_piece", "nanatsu_no_taizai", "shokugeki_no_soma", "fairy_tail", "boruto"}

	for _, tag := range popularTags {
		if titleId == tag {
			return true
		}
	}
	return false
}

func (s *mangaService) GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error) {
	ml, err := s.mClient.GetMangaList()
	if err != nil {
		return nil, nil, mErr.NewGenericError()
	}

	if len(ml.Mangas) == 0 {
		return nil, nil, mErr.NewNotFoundError("manga")
	}

	var pMangas []contract.Manga
	var lMangas []contract.Manga
	for _, dm := range ml.Mangas {
		manga := getMappedManga(dm)
		if isPopularManga(manga.TitleId) {
			pMangas = append(pMangas, manga)
		} else {
			lMangas = append(lMangas, manga)
		}
	}

	if len(pMangas) == 0 {
		return nil, &lMangas, nil
	}

	if len(lMangas) == 0 {
		return &pMangas, nil, nil
	}

	return &pMangas, &lMangas, nil
}

func NewMangaService(mClient client.MangaClient) *mangaService {
	return &mangaService{mClient: mClient}
}
