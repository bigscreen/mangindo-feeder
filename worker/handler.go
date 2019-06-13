package worker

import (
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
)

func InitWorkerHandler(w adapter.Worker, d service.WorkerDependencies) {
	registerSetMangaCacheJob(w, d)
	registerSetChapterCacheJob(w, d)
	registerSetContentCacheJob(w, d)
}

func registerSetMangaCacheJob(w adapter.Worker, d service.WorkerDependencies) {
	err := w.Register(constants.SetMangaCacheJob, func(args adapter.Args) error {
		return d.MangaCacheManager.SetCache()
	})
	if err != nil {
		logger.Errorf("Error while registering %s job, error: %s", constants.SetMangaCacheJob, err.Error())
	}
}

func registerSetChapterCacheJob(w adapter.Worker, d service.WorkerDependencies) {
	err := w.Register(constants.SetChapterCacheJob, func(args adapter.Args) error {
		titleId := args[constants.JobArgTitleId].(string)
		return d.ChapterCacheManager.SetCache(titleId)
	})
	if err != nil {
		logger.Errorf("Error while registering %s job, error: %s", constants.SetChapterCacheJob, err.Error())
	}
}

func registerSetContentCacheJob(w adapter.Worker, d service.WorkerDependencies) {
	err := w.Register(constants.SetContentCacheJob, func(args adapter.Args) error {
		titleId := args[constants.JobArgTitleId].(string)
		chapter := args[constants.JobArgChapter].(float64)
		return d.ContentCacheManager.SetCache(titleId, float32(chapter))
	})
	if err != nil {
		logger.Errorf("Error while registering %s job, error: %s", constants.SetContentCacheJob, err.Error())
	}
}
