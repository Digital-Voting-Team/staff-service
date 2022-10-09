package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetPositionRequest struct {
	PositionID int64 `url:"-"`
}

func NewGetPositionRequest(r *http.Request) (GetPositionRequest, error) {
	request := GetPositionRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.PositionID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
