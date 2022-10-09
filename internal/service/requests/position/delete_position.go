package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeletePositionRequest struct {
	PositionID int64 `url:"-"`
}

func NewDeletePositionRequest(r *http.Request) (DeletePositionRequest, error) {
	request := DeletePositionRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.PositionID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
