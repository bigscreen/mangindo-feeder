package service

import (
	"errors"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/logger"
	mMock "github.com/bigscreen/mangindo-feeder/mock"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WorkerServiceTestSuite struct {
	suite.Suite
}

func TestWorkerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerServiceTestSuite))
}

func (s *WorkerServiceTestSuite) SetupSuite() {
	config.Load()
	appcontext.Initiate()
	logger.SetupLogger()
}

func (s *WorkerServiceTestSuite) TestSetMangaCache_ReturnsNil_WhenItSucceeds() {
	w := &mMock.WorkerAdapterMock{}
	stubSetMangaJob(w, nil)

	ws := NewWorkerService(w)
	err := ws.SetMangaCache()
	assert.Nil(s.T(), err)
	w.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestSetMangaCache_ReturnsError_WhenItFails() {
	w := &mMock.WorkerAdapterMock{}
	stubSetMangaJob(w, errors.New("some error"))

	ws := NewWorkerService(w)
	err := ws.SetMangaCache()
	assert.NotNil(s.T(), err)
	w.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestSetChapterCache_ReturnsNil_WhenItSucceeds() {
	w := &mMock.WorkerAdapterMock{}
	stubSetChapterJob(w, "bleach", nil)

	ws := NewWorkerService(w)
	err := ws.SetChapterCache("bleach")
	assert.Nil(s.T(), err)
	w.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestSetChapterCache_ReturnsError_WhenItFails() {
	w := &mMock.WorkerAdapterMock{}
	stubSetChapterJob(w, "bleach", errors.New("some error"))

	ws := NewWorkerService(w)
	err := ws.SetChapterCache("bleach")
	assert.NotNil(s.T(), err)
	w.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestSetContentCache_ReturnsNil_WhenItSucceeds() {
	w := &mMock.WorkerAdapterMock{}
	stubSetContentJob(w, "bleach", float32(650), nil)

	ws := NewWorkerService(w)
	err := ws.SetContentCache("bleach", float32(650))
	assert.Nil(s.T(), err)
	w.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestSetContentCache_ReturnsError_WhenItFails() {
	w := &mMock.WorkerAdapterMock{}
	stubSetContentJob(w, "bleach", float32(650), errors.New("some error"))

	ws := NewWorkerService(w)
	err := ws.SetContentCache("bleach", float32(650))
	assert.NotNil(s.T(), err)
	w.AssertExpectations(s.T())
}

func stubSetMangaJob(w *mMock.WorkerAdapterMock, returnedErr error) {
	stubWorkerPerform(w, constants.SetMangaCacheJob, nil).Return(returnedErr)
}

func stubSetChapterJob(w *mMock.WorkerAdapterMock, titleId string, returnedErr error) {
	stubWorkerPerform(w, constants.SetChapterCacheJob, adapter.Args{
		constants.JobArgTitleId: titleId,
	}).Return(returnedErr)
}

func stubSetContentJob(w *mMock.WorkerAdapterMock, titleId string, chapter float32, returnedErr error) {
	stubWorkerPerform(w, constants.SetContentCacheJob, adapter.Args{
		constants.JobArgTitleId: titleId,
		constants.JobArgChapter: chapter,
	}).Return(returnedErr)
}

func stubWorkerPerform(w *mMock.WorkerAdapterMock, handlerName string, args adapter.Args) *mock.Call {
	return w.On("Perform", adapter.Job{
		Queue:   constants.WorkerDefaultQueue,
		Handler: handlerName,
		Args:    args,
	})
}
