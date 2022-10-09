package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetStaffRequest struct {
	StaffID int64 `url:"-"`
}

func NewGetStaffRequest(r *http.Request) (GetStaffRequest, error) {
	request := GetStaffRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.StaffID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
