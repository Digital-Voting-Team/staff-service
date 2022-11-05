package handlers

import (
	"github.com/Digital-Voting-Team/staff-service/internal/data"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/position"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetPositionList(w http.ResponseWriter, r *http.Request) {
	accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	if *accessLevel < resources.Manager {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewGetPositionListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	positionsQ := helpers.PositionsQ(r)
	applyFilters(positionsQ, request)
	position, err := positionsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get position")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.PositionListResponse{
		Data:  newPositionesList(position),
		Links: helpers.GetOffsetLinks(r, request.OffsetPageParams),
	}
	ape.Render(w, response)
}

func applyFilters(q data.PositionsQ, request requests.GetPositionListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterName) > 0 {
		q.FilterByName(request.FilterName...)
	}

	if len(request.FilterAccessLevel) > 0 {
		q.FilterByAccessLevel(request.FilterAccessLevel...)
	}
}

func newPositionesList(positiones []data.Position) []resources.Position {
	result := make([]resources.Position, len(positiones))
	for i, position := range positiones {
		result[i] = resources.Position{
			Key: resources.NewKeyInt64(position.ID, resources.POSITION),
			Attributes: resources.PositionAttributes{
				Name:        position.Name,
				AccessLevel: *position.AccessLevel,
			},
		}
	}
	return result
}
