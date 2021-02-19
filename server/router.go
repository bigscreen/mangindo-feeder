package server

import (
	"net/http"

	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/handler"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/gorilla/mux"
)

func Router(deps service.Dependencies) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	router.HandleFunc(constants.GetMangasAPIPath, handler.GetMangas(deps.MangaService)).Methods("GET")
	router.HandleFunc(constants.GetChaptersAPIPath, handler.GetChapters(deps.ChapterService)).Methods("GET")
	router.HandleFunc(constants.GetContentsAPIPath, handler.GetContents(deps.ContentService)).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	return router
}
