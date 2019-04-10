package adapter

import (
	"context"

	"fmt"

	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type Options struct {
	Pool           *redis.Pool
	Name           string
	MaxConcurrency int
}

var _ Worker = &Adapter{}

type Adapter struct {
	Enqueur *work.Enqueuer
	Pool    *work.WorkerPool
	ctx     context.Context
	Client  *work.Client
	Name    string
}

func NewAdapter(opts Options) *Adapter {
	ctx := context.Background()
	if opts.Name == "" {
		opts.Name = "goWorker"
	}
	if opts.MaxConcurrency == 0 {
		opts.MaxConcurrency = 25
	}
	enqueuer := work.NewEnqueuer(opts.Name, opts.Pool)
	client := work.NewClient(opts.Name, opts.Pool)
	pool := work.NewWorkerPool(struct{}{}, uint(opts.MaxConcurrency), opts.Name, opts.Pool)
	client.Queues()
	return &Adapter{
		Enqueur: enqueuer,
		Pool:    pool,
		ctx:     ctx,
		Name:    opts.Name,
		Client:  client,
	}
}

func (q *Adapter) Start(ctx context.Context) error {
	fmt.Println("Starting gocraft/work Worker")
	q.ctx = ctx
	go func() {
		select {
		case <-ctx.Done():
			q.Stop()
		}
	}()
	q.Pool.Start()
	return nil
}
func (q *Adapter) Stop() error {
	fmt.Println("Stopping gocraft/work Worker")
	q.Pool.Stop()
	return nil
}

func (q *Adapter) Register(name string, h Handler) error {
	q.Pool.Job(name, func(job *work.Job) error {
		return h(job.Args)
	})
	return nil
}

func (q *Adapter) RegisterWithRetrial(name string, h Handler, retry uint) error {
	opts := work.JobOptions{}
	opts.MaxFails = retry
	q.Pool.JobWithOptions(name, opts, func(job *work.Job) error {
		return h(job.Args)
	})
	return nil
}

func (q *Adapter) Perform(job Job) error {
	fmt.Printf("Enqueuing job %s\n", job)

	_, err := q.Enqueur.Enqueue(job.Handler, job.Args)
	if err != nil {
		fmt.Printf("error enqueuing job %s", job)
		return errors.WithStack(err)
	}
	return nil
}

func (q *Adapter) PerformUnique(job Job) error {
	fmt.Printf("Enqueuing unique job %s\n", job)

	_, err := q.Enqueur.EnqueueUnique(job.Handler, job.Args)
	if err != nil {
		fmt.Printf("Error enqueuing unique job %s", job)
		return errors.WithStack(err)
	}
	return nil
}

func (q *Adapter) PerformIn(job Job, t time.Duration) error {
	fmt.Printf("Enqueuing job %s\n", job)
	d := int64(t / time.Second)

	_, err := q.Enqueur.EnqueueIn(job.Handler, d, job.Args)
	if err != nil {
		fmt.Printf("error enqueuing job %s", job)
		return errors.WithStack(err)
	}
	return nil
}

func (q *Adapter) PerformPeriodically(cronSchedule string, job Job) error {
	fmt.Printf("Enqueuing job %s\n", job)

	q.Pool.PeriodicallyEnqueue(cronSchedule, job.Handler)
	return nil
}

func (q *Adapter) PerformAt(job Job, t time.Time) error {
	return q.PerformIn(job, t.Sub(time.Now()))
}

func Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}
