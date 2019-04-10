package main

import (
	"github.com/bigscreen/mangindo-feeder/appcontext"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/server"
	"github.com/bigscreen/mangindo-feeder/worker"
	"github.com/urfave/cli"
	"os"
)

func handleInitError() {
	if e := recover(); e != nil {
		logger.Fatalf("Failed to load the app due to error : %s", e)
	}
}

func main() {
	defer handleInitError()

	config.Load()
	logger.SetupLogger()

	appcontext.Initiate()

	clientApp := cli.NewApp()
	clientApp.Name = "mangindo-feeder"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "start",
			Description: "Start HTTP api server",
			Action: func(c *cli.Context) error {
				server.StartAPIServer()
				return nil
			},
		}, {
			Name:        "worker",
			Description: "Start worker process",
			Action: func(c *cli.Context) error {
				worker.Start()
				return nil
			},
		}, {
			Name:        "web-worker",
			Description: "Start worker web process",
			Action: func(c *cli.Context) error {
				worker.StartWorkerWebServer()
				return nil
			},
		},
	}

	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
}
