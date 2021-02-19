package service

import (
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/common"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
)

type ChapterService interface {
	GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error)
}

type chapterService struct {
	chapterClient       client.ChapterClient
	chapterCacheManager manager.ChapterCacheManager
	workerService       WorkerService
}

func (s *chapterService) GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error) {
	cl, err := s.chapterCacheManager.GetCache(req.TitleID)
	if err != nil {
		cl, err = s.chapterClient.GetChapterList(req.TitleID)
		if err != nil {
			return nil, mErr.NewGenericError()
		}

		err = s.workerService.SetChapterCache(req.TitleID)
		if err != nil {
			logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetChapterCacheJob, err.Error())
		}
	}

	if len(cl.Chapters) == 0 {
		return nil, mErr.NewNotFoundError("chapter")
	}

	var chapters []contract.Chapter
	for _, dc := range cl.Chapters {
		chapter := contract.Chapter{
			Number:  common.GetFormattedChapterNumber(dc.Number),
			Title:   dc.Title,
			TitleID: dc.TitleID,
		}
		chapters = append(chapters, chapter)
	}

	return &chapters, nil
}

func NewChapterService(cc client.ChapterClient, ccm manager.ChapterCacheManager, ws WorkerService) *chapterService {
	return &chapterService{
		chapterClient:       cc,
		chapterCacheManager: ccm,
		workerService:       ws,
	}
}
