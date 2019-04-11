package server

import (
	"github.com/bigscreen/mangindo-feeder/handler"
	"github.com/gorilla/mux"
	"net/http"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	return router
}
