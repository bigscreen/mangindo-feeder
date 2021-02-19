package service

import (
	"github.com/bigscreen/mangindo-feeder/constants"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
)

type workerService struct {
	adapter adapter.Worker
}

type WorkerService interface {
	SetMangaCache() error
	SetChapterCache(titleID string) error
	SetContentCache(titleID string, chapter float32) error
}

func (s *workerService) SetMangaCache() error {
	err := s.adapter.Perform(adapter.Job{
		Queue:   constants.WorkerDefaultQueue,
		Handler: constants.SetMangaCacheJob,
	})
	if err != nil {
		logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetMangaCacheJob, err.Error())
		return mErr.NewWorkerError(err.Error())
	}

	return nil
}

func (s *workerService) SetChapterCache(titleID string) error {
	err := s.adapter.Perform(adapter.Job{
		Queue:   constants.WorkerDefaultQueue,
		Handler: constants.SetChapterCacheJob,
		Args: adapter.Args{
			constants.JobArgTitleID: titleID,
		},
	})
	if err != nil {
		logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetChapterCacheJob, err.Error())
		return mErr.NewWorkerError(err.Error())
	}

	return nil
}

func (s *workerService) SetContentCache(titleID string, chapter float32) error {
	err := s.adapter.Perform(adapter.Job{
		Queue:   constants.WorkerDefaultQueue,
		Handler: constants.SetContentCacheJob,
		Args: adapter.Args{
			constants.JobArgTitleID: titleID,
			constants.JobArgChapter: chapter,
		},
	})
	if err != nil {
		logger.Errorf("Failed to enqueue %s job, with error: %s", constants.SetContentCacheJob, err.Error())
		return mErr.NewWorkerError(err.Error())
	}

	return nil
}

func NewWorkerService(adapter adapter.Worker) *workerService {
	return &workerService{
		adapter: adapter,
	}
}
