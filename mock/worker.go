package mock

import (
	"context"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
	"github.com/stretchr/testify/mock"
	"time"
)

type WorkerAdapterMock struct {
	mock.Mock
}

func (m WorkerAdapterMock) Start(c context.Context) error {
	args := m.Called(c)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) Stop() error {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) Perform(job adapter.Job) error {
	args := m.Called(job)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) PerformUnique(job adapter.Job) error {
	args := m.Called(job)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) Register(s string, handler adapter.Handler) error {
	args := m.Called(s, handler)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) RegisterWithRetrial(s string, handler adapter.Handler, retrial uint) error {
	args := m.Called(s, handler, retrial)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) PerformIn(job adapter.Job, t time.Duration) error {
	args := m.Called(job, t)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) PerformAt(job adapter.Job, t time.Time) error {
	args := m.Called(job, t)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m WorkerAdapterMock) PerformPeriodically(s string, job adapter.Job) error {
	args := m.Called(s, job)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}
