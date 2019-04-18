package service

import "github.com/bigscreen/mangindo-feeder/client"

type Dependencies struct {
	MangaService   MangaService
	ChapterService ChapterService
	ContentService ContentService
}

func InstantiateDependencies() Dependencies {
	mangaClient := client.NewMangaClient()
	chapterClient := client.NewChapterClient()
	contentClient := client.NewContentClient()

	mangaService := NewMangaService(mangaClient)
	chapterService := NewChapterService(chapterClient)
	contentService := NewContentService(contentClient)

	return Dependencies{
		MangaService:   mangaService,
		ChapterService: chapterService,
		ContentService: contentService,
	}
}
