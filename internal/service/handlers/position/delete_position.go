package handlers

import (
	"net/http"
	"staff-service/internal/service/helpers"
	requests "staff-service/internal/service/requests/position"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeletePosition(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeletePositionRequest(r)
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

	err = helpers.PositionsQ(r).Delete(request.PositionID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete position")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
