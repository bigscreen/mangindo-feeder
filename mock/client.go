package mock

import (
	"github.com/bigscreen/mangindo-feeder/domain"
	"github.com/stretchr/testify/mock"
)

type MangaClientMock struct {
	mock.Mock
}

func (m *MangaClientMock) GetMangaList() (*domain.MangaListResponse, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.MangaListResponse), nil
}

type ChapterClientMock struct {
	mock.Mock
}

func (m *ChapterClientMock) GetChapterList(titleID string) (*domain.ChapterListResponse, error) {
	args := m.Called(titleID)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.ChapterListResponse), nil
}

type ContentClientMock struct {
	mock.Mock
}

func (m *ContentClientMock) GetContentList(titleID string, chapter float32) (*domain.ContentListResponse, error) {
	args := m.Called(titleID, chapter)
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*domain.ContentListResponse), nil
}
