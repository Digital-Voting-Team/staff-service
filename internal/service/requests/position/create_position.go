package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreatePositionRequest struct {
	Data resources.Position
}

func NewCreatePositionRequest(r *http.Request) (CreatePositionRequest, error) {
	var request CreatePositionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreatePositionRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/name": validation.Validate(&r.Data.Attributes.Name, validation.Required,
			validation.Length(3, 45)),
		"/data/attributes/access_level": validation.Validate(&r.Data.Attributes.AccessLevel, validation.Required,
			validation.By(helpers.IsValidAccessLevel)),
	}).Filter()
}
