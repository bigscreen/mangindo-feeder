package server

import (
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/handler"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(deps service.Dependencies) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	router.HandleFunc(constants.GetMangasApiPath, handler.GetMangas(deps.MangaService)).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	return router
}
