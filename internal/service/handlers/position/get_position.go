package handlers

import (
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/position"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetPosition(w http.ResponseWriter, r *http.Request) {
	accessLevel := r.Context().Value("accessLevel").(resources.AccessLevel)
	if accessLevel < resources.Manager {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewGetPositionRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resultPosition, err := helpers.PositionsQ(r).FilterByID(request.PositionID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get resultPosition from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if resultPosition == nil {
		ape.Render(w, problems.NotFound())
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
