package service

import (
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/bigscreen/mangindo-feeder/domain"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
)

type MangaService interface {
	GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error)
}

type mangaService struct {
	mangaClient       client.MangaClient
	mangaCacheManager manager.MangaCacheManager
	workerService     WorkerService
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
	for _, tag := range config.PopularMangaTags() {
		if titleId == tag {
			return true
		}
	}
	return false
}

func (s *mangaService) GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error) {
	ml, err := s.mangaCacheManager.GetCache()
	if err != nil {
		ml, err = s.mangaClient.GetMangaList()
		if err != nil {
			return nil, nil, mErr.NewGenericError()
		}

		err = s.workerService.SetMangaCache()
		if err != nil {
			logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetMangaCacheJob, err.Error())
		}
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

func NewMangaService(mc client.MangaClient, mcm manager.MangaCacheManager, ws WorkerService) *mangaService {
	return &mangaService{
		mangaClient:       mc,
		mangaCacheManager: mcm,
		workerService:     ws,
	}
}
