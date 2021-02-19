package worker

import (
	"fmt"

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
		titleID, ok := args[constants.JobArgTitleID].(string)
		if !ok {
			return fmt.Errorf("can not get argument %s", constants.JobArgTitleID)
		}
		return d.ChapterCacheManager.SetCache(titleID)
	})
	if err != nil {
		logger.Errorf("Error while registering %s job, error: %s", constants.SetChapterCacheJob, err.Error())
	}
}

func registerSetContentCacheJob(w adapter.Worker, d service.WorkerDependencies) {
	err := w.Register(constants.SetContentCacheJob, func(args adapter.Args) error {
		titleID, ok := args[constants.JobArgTitleID].(string)
		if !ok {
			return fmt.Errorf("can not get argument %s", constants.JobArgTitleID)
		}
		chapter, ok := args[constants.JobArgChapter].(float64)
		if !ok {
			return fmt.Errorf("can not get argument %s", constants.JobArgChapter)
		}
		return d.ContentCacheManager.SetCache(titleID, float32(chapter))
	})
	if err != nil {
		logger.Errorf("Error while registering %s job, error: %s", constants.SetContentCacheJob, err.Error())
	}
}
