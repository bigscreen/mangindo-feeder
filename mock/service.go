package mock

import (
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/stretchr/testify/mock"
)

type MockMangaService struct {
	mock.Mock
}

func (m MockMangaService) GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error) {
	args := m.Called()
	if args.Get(2) != nil {
		return nil, nil, args.Get(2).(error)
	}
	if args.Get(0) != nil && args.Get(1) == nil {
		return args.Get(0).(*[]contract.Manga), nil, nil
	}
	if args.Get(1) != nil && args.Get(0) == nil {
		return nil, args.Get(1).(*[]contract.Manga), nil
	}
	return args.Get(0).(*[]contract.Manga), args.Get(1).(*[]contract.Manga), nil
}

type MockChapterService struct {
	mock.Mock
}

func (m MockChapterService) GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error) {
	args := m.Called(req)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*[]contract.Chapter), nil
}

type MockContentService struct {
	mock.Mock
}

func (m MockContentService) GetContents(req contract.ContentRequest) (*[]contract.Content, error) {
	args := m.Called(req)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*[]contract.Content), nil
}
