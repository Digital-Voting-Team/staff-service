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

func UpdatePosition(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdatePositionRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	position, err := helpers.PositionsQ(r).FilterByID(request.PositionID).Get()
	if position == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	newPosition := data.Position{
		Name:        request.Data.Attributes.Name,
		AccessLevel: &request.Data.Attributes.AccessLevel,
	}

	var resultPosition data.Position
	resultPosition, err = helpers.PositionsQ(r).FilterByID(position.ID).Update(newPosition)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update position")
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
