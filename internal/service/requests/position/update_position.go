package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"staff-service/internal/service/helpers"
	"staff-service/resources"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdatePositionRequest struct {
	PositionID int64 `url:"-" json:"-"`
	Data       resources.Position
}

func NewUpdatePositionRequest(r *http.Request) (UpdatePositionRequest, error) {
	request := UpdatePositionRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.PositionID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdatePositionRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/name": validation.Validate(&r.Data.Attributes.Name, validation.Required,
			validation.Length(3, 45)),
		"/data/attributes/access_level": validation.Validate(&r.Data.Attributes.AccessLevel, validation.Required,
			validation.By(helpers.IsValidAccessLevel)),
	}).Filter()
}
