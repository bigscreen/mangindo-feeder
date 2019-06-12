package service

import (
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/common"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
)

type ChapterService interface {
	GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error)
}

type chapterService struct {
	cClient client.ChapterClient
}

func (s *chapterService) GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error) {
	cl, err := s.cClient.GetChapterList(req.TitleId)
	if err != nil {
		return nil, mErr.NewGenericError()
	}

	if len(cl.Chapters) == 0 {
		return nil, mErr.NewNotFoundError("chapter")
	}

	var chapters []contract.Chapter
	for _, dc := range cl.Chapters {
		chapter := contract.Chapter{
			Number:  common.GetFormattedChapterNumber(dc.Number),
			Title:   dc.Title,
			TitleId: dc.TitleId,
		}
		chapters = append(chapters, chapter)
	}

	return &chapters, nil
}

func NewChapterService(cClient client.ChapterClient) *chapterService {
	return &chapterService{cClient: cClient}
}
