package client

import (
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/stretchr/testify/mock"
)

type MockMangaClient struct {
	mock.Mock
}

func (m MockMangaClient) GetMangaList() (*domain.MangaListResponse, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.MangaListResponse), nil
}

type MockChapterClient struct {
	mock.Mock
}

func (m MockChapterClient) GetChapterList(titleId string) (*domain.ChapterListResponse, error) {
	args := m.Called(titleId)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.ChapterListResponse), nil
}

type MockContentClient struct {
	mock.Mock
}

func (m MockContentClient) GetContentList(titleId string, chapter float32) (*domain.ContentListResponse, error) {
	args := m.Called(titleId, chapter)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.ContentListResponse), nil
}
