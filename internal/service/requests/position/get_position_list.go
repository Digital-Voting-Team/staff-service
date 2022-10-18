package requests

import (
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetPositionListRequest struct {
	pgdb.OffsetPageParams
	FilterName        []string                `filter:"name"`
	FilterAccessLevel []resources.AccessLevel `filter:"access_level"`
}

func NewGetPositionListRequest(r *http.Request) (GetPositionListRequest, error) {
	var request GetPositionListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
