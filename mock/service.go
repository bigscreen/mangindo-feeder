package mock

import (
	"github.com/bigscreen/mangindo-feeder/contract"
	"github.com/stretchr/testify/mock"
)

type MangaServiceMock struct {
	mock.Mock
}

func (m MangaServiceMock) GetMangas() (popular *[]contract.Manga, latest *[]contract.Manga, err error) {
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

type ChapterServiceMock struct {
	mock.Mock
}

func (m ChapterServiceMock) GetChapters(req contract.ChapterRequest) (*[]contract.Chapter, error) {
	args := m.Called(req)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*[]contract.Chapter), nil
}

type ContentServiceMock struct {
	mock.Mock
}

func (m ContentServiceMock) GetContents(req contract.ContentRequest) (*[]contract.Content, error) {
	args := m.Called(req)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*[]contract.Content), nil
}

type WorkerServiceMock struct {
	mock.Mock
}

func (m WorkerServiceMock) SetMangaCache() error {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerServiceMock) SetChapterCache(titleId string) error {
	args := m.Called(titleId)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerServiceMock) SetContentCache(titleId string, chapter float32) error {
	args := m.Called(titleId, chapter)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}
