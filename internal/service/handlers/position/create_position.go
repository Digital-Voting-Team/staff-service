package handlers

import (
	"net/http"
	"staff-service/internal/data"
	"staff-service/internal/service/helpers"
	requests "staff-service/internal/service/requests/position"
	"staff-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreatePosition(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreatePositionRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultPosition data.Position

	position := data.Position{
		Name:        request.Data.Attributes.Name,
		AccessLevel: &request.Data.Attributes.AccessLevel,
	}

	resultPosition, err = helpers.PositionsQ(r).Insert(position)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create position")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.PositionResponse{
		Data: resources.Position{
			Key: resources.NewKeyInt64(resultPosition.ID, resources.POSITION),
			Attributes: resources.PositionAttributes{
				Name:        resultPosition.Name,
				AccessLevel: *resultPosition.AccessLevel,
			},
		},
	}
	ape.Render(w, result)
}
