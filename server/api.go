package server

import (
	"context"
	"github.com/bigscreen/mangindo-feeder/config"
	"github.com/bigscreen/mangindo-feeder/logger"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/codegangsta/negroni"
	"github.com/getsentry/raven-go"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func StartAPIServer() {
	logger.Info("Starting mangindo-feeder service")

	deps := service.InstantiateDependencies()
	muxRouter := Router(deps)
	handlerFunc := muxRouter.ServeHTTP

	n := negroni.New(negroni.NewRecovery())
	n.Use(negroniRecoverHandler())
	n.UseHandlerFunc(handlerFunc)
	portInfo := ":" + strconv.Itoa(config.Port())
	server := &http.Server{Addr: portInfo, Handler: n}
	go listenServer(server)
	waitForShutdown(server)
}

func negroniRecoverHandler() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fn := raven.RecoveryHandler(next)
		fn(w, r)
	})
}

func listenServer(apiServer *http.Server) {
	err := apiServer.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Fatal(err.Error())
	}
}

func waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Info("Mangindo-feeder is shutting down")
	apiServer.Shutdown(context.Background())
	logger.Info("Mangindo-feeder shutdown complete")
}
