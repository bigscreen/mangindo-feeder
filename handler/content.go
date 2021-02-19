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

func GetContents(s service.ContentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		titleID := vars[constants.TitleIDKeyParam]
		chapter := vars[constants.ChapterKeyParam]

		validators := []validator.Validator{
			validator.PresenceValidator{Field: constants.TitleIDKeyParam, Value: &titleID},
			validator.PresenceValidator{Field: constants.ChapterKeyParam, Value: &chapter},
			validator.NumberValidator{Field: constants.ChapterKeyParam, Value: &chapter},
		}
		isValid, errMsgs := validator.ValidateAll(validators)
		if !isValid {
			err := mErr.NewValidationError(errMsgs)
			respondWith(mErr.GetStatusCodeOf(err), r, w, getErrorResponse(err))
			return
		}

		contents, err := s.GetContents(contract.NewContentRequest(titleID, chapter))
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
