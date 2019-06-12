package worker

import (
	"context"
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/bigscreen/mangindo-feeder/worker/adapter"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	wd := service.InstantiateWorkerDependencies()
	wa := appcontext.GetWorkerAdapter()
	InitWorkerHandler(wa, wd)
	err := wa.Start(context.Background())
	if err != nil {
		logger.Error("Failed to start worker")
	}
	logger.Info("Starting worker")

	waitForWorkerShutDown(wa)
}

func waitForWorkerShutDown(worker adapter.Worker) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	reason := <-sig
	logger.Info("Worker is shutting down ", reason)
	err := worker.Stop()
	if err != nil {
		logger.Error("Failed to shut down worker")
	}
	logger.Info("Worker shutdown complete")
}

func StartWorkerWebServer() {
	logger.Debug("Starting worker web server")
}
