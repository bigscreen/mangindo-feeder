package adapter

import (
	"context"
	"time"
)

type Handler func(Args) error

type Worker interface {
	Start(context.Context) error
	Stop() error
	Perform(Job) error
	PerformUnique(Job) error
	Register(string, Handler) error
	RegisterWithRetrial(string, Handler, uint) error
	PerformIn(Job, time.Duration) error
	PerformAt(Job, time.Time) error
	PerformPeriodically(string, Job) error
}
