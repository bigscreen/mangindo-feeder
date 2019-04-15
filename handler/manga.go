package handler

import (
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/service"
	"net/http"
)

func GetMangas(s service.MangaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pop, lts, err := s.GetMangas()
		if err != nil {
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		var pms []contract.Manga
		if pop != nil {
			pms = *pop
		} else {
			pms = []contract.Manga{}
		}

		var lms []contract.Manga
		if lts != nil {
			lms = *lts
		} else {
			lms = []contract.Manga{}
		}

		mr := contract.MangaResponse{
			Success:       true,
			PopularMangas: pms,
			LatestMangas:  lms,
		}
		respondWith(http.StatusOK, r, w, mr)

	}
}
