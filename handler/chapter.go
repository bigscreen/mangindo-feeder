package handler

import (
	"net/http"

	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/bigscreen/mangindo-feeder/validator"
	"github.com/gorilla/mux"
)

func GetChapters(s service.ChapterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		titleID := vars[constants.TitleIDKeyParam]

		validators := []validator.Validator{
			validator.PresenceValidator{Field: constants.TitleIDKeyParam, Value: &titleID},
		}
		isValid, errMsgs := validator.ValidateAll(validators)
		if !isValid {
			err := mErr.NewValidationError(errMsgs)
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		chapters, err := s.GetChapters(contract.NewChapterRequest(titleID))
		if err != nil {
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		cr := contract.ChapterResponse{
			Success:  true,
			Chapters: *chapters,
		}
		respondWith(http.StatusOK, r, w, cr)
	}
}
