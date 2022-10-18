package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/staff-service/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateStaffRequest struct {
	StaffID int64 `url:"-" json:"-"`
	Data    resources.Staff
}

func NewUpdateStaffRequest(r *http.Request) (UpdateStaffRequest, error) {
	request := UpdateStaffRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.StaffID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateStaffRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/employment_date": validation.Validate(&r.Data.Attributes.EmploymentDate,
			validation.Required, validation.By(helpers.IsDate)),
		"/data/attributes/salary": validation.Validate(&r.Data.Attributes.Salary,
			validation.Required),
		"/data/attributes/status": validation.Validate(&r.Data.Attributes.Status,
			validation.Required, validation.By(helpers.IsValidWorkerStatus)),
		"/data/relationships/person/data/id": validation.Validate(&r.Data.Relationships.Person.Data.ID,
			validation.Required, validation.By(helpers.IsInteger)),
		"/data/relationships/person/cafe/id": validation.Validate(&r.Data.Relationships.Cafe.Data.ID,
			validation.Required, validation.By(helpers.IsInteger)),
		"/data/relationships/person/position/id": validation.Validate(&r.Data.Relationships.Position.Data.ID,
			validation.Required, validation.By(helpers.IsInteger)),
	}).Filter()
}
