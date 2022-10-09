package requests

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetStaffListRequest struct {
	pgdb.OffsetPageParams
	FilterWorkStart    *time.Time `filter:"work_start"`
	FilterWorkEnd      *time.Time `filter:"work_end"`
	FilterSalaryUp     []float32  `filter:"salary_lowest"`
	FilterSalaryBottom []float32  `filter:"salary_greatest"`
	FilterPosition     []int64    `filter:"position_id"`
	FilterCafe         []int64    `filter:"cafe_id"`
}

func NewGetStaffListRequest(r *http.Request) (GetStaffListRequest, error) {
	var request GetStaffListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
