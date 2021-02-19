package manager

import (
	"encoding/json"
	"errors"

	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/client"
	"github.com/bigscreen/mangindo-feeder/domain"
)

type mangaCacheManager struct {
	mClient client.MangaClient
	mCache  cache.MangaCache
}

type MangaCacheManager interface {
	SetCache() error
	GetCache() (*domain.MangaListResponse, error)
}

func (m *mangaCacheManager) SetCache() error {
	ml, err := m.mClient.GetMangaList()
	if err != nil {
		return err
	}

	ms, _ := json.Marshal(ml)

	return m.mCache.Set(string(ms))
}

func (m *mangaCacheManager) GetCache() (*domain.MangaListResponse, error) {
	ms, err := m.mCache.Get()
	if err != nil {
		return nil, err
	}

	var ml *domain.MangaListResponse
	err = json.Unmarshal([]byte(ms), &ml)
	if err != nil {
		return nil, errors.New("invalid manga cache")
	}

	return ml, nil
}

func NewMangaCacheManager(client client.MangaClient, cache cache.MangaCache) *mangaCacheManager {
	return &mangaCacheManager{
		mClient: client,
		mCache:  cache,
	}
}
