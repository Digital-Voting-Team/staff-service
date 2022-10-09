package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteStaffRequest struct {
	StaffID int64 `url:"-"`
}

func NewDeleteStaffRequest(r *http.Request) (DeleteStaffRequest, error) {
	request := DeleteStaffRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.StaffID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
