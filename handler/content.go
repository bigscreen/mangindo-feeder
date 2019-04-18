package handler

import (
	"github.com/bigscreen/mangindo-feeder/constants"
	"github.com/bigscreen/mangindo-feeder/contract"
	mErr "github.com/bigscreen/mangindo-feeder/error"
	"github.com/bigscreen/mangindo-feeder/service"
	"github.com/bigscreen/mangindo-feeder/validator"
	"net/http"
)

func GetContents(s service.ContentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		titleId := r.URL.Query().Get(constants.TitleIdKeyParam)
		chapter := r.URL.Query().Get(constants.ChapterKeyParam)

		validators := []validator.Validator{
			validator.PresenceValidator{Field: constants.TitleIdKeyParam, Value: &titleId},
			validator.PresenceValidator{Field: constants.ChapterKeyParam, Value: &chapter},
			validator.NumberValidator{Field: constants.ChapterKeyParam, Value: &chapter},
		}
		isValid, errMsgs := validator.ValidateAll(validators)
		if !isValid {
			err := mErr.NewValidationError(errMsgs)
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		contents, err := s.GetContents(contract.NewContentRequest(titleId, chapter))
		if err != nil {
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		cr := contract.ContentResponse{
			Success:  true,
			Contents: *contents,
		}
		respondWith(http.StatusOK, r, w, cr)
	}
}
