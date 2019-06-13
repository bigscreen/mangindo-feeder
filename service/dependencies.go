package service

import (
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/cache"
	"github.com/bigscreen/mangindo-feeder/cache/manager"
	"github.com/bigscreen/mangindo-feeder/client"
)

type Dependencies struct {
	MangaService   MangaService
	ChapterService ChapterService
	ContentService ContentService
}

type WorkerDependencies struct {
	MangaCacheManager   manager.MangaCacheManager
	ChapterCacheManager manager.ChapterCacheManager
	ContentCacheManager manager.ContentCacheManager
}

func InstantiateDependencies() Dependencies {
	macl := client.NewMangaClient()
	chcl := client.NewChapterClient()
	cocl := client.NewContentClient()

	maca := cache.NewMangaCache()
	chca := cache.NewChapterCache()
	coca := cache.NewContentCache()

	macm := manager.NewMangaCacheManager(macl, maca)
	chcm := manager.NewChapterCacheManager(chcl, chca)
	cocm := manager.NewContentCacheManager(cocl, coca)

	ws := NewWorkerService(appcontext.GetWorkerAdapter())

	mas := NewMangaService(macl, macm, ws)
	chs := NewChapterService(chcl, chcm, ws)
	cos := NewContentService(cocl, cocm, ws)

	return Dependencies{
		MangaService:   mas,
		ChapterService: chs,
		ContentService: cos,
	}
}

func InstantiateWorkerDependencies() WorkerDependencies {
	mangaClient := client.NewMangaClient()
	chapterClient := client.NewChapterClient()
	contentClient := client.NewContentClient()

	mangaCache := cache.NewMangaCache()
	chapterCache := cache.NewChapterCache()
	contentCache := cache.NewContentCache()

	mangaCacheManager := manager.NewMangaCacheManager(mangaClient, mangaCache)
	chapterCacheManager := manager.NewChapterCacheManager(chapterClient, chapterCache)
	contentCacheManager := manager.NewContentCacheManager(contentClient, contentCache)

	return WorkerDependencies{
		MangaCacheManager:   mangaCacheManager,
		ChapterCacheManager: chapterCacheManager,
		ContentCacheManager: contentCacheManager,
	}
}
