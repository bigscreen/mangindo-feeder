package handler

import (
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/bigscreen/mangindo-feeder/validator"
	"net/http"
)

func GetChapters(s service.ChapterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		titleId := r.URL.Query().Get(constants.TitleIdKeyParam)

		validators := []validator.Validator{
			validator.PresenceValidator{Field: constants.TitleIdKeyParam, Value: &titleId},
		}
		isValid, errMsgs := validator.ValidateAll(validators)
		if !isValid {
			err := mErr.NewValidationError(errMsgs)
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		chapters, err := s.GetChapters(contract.NewChapterRequest(titleId))
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
