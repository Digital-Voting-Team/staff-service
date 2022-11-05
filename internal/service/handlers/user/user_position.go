package user

import (
	"errors"
	"github.com/Digital-Voting-Team/staff-service/internal/data"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/user"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetPositionByUser(r *http.Request, userId int64) (*data.Position, error) {
	resultStaff, err := helpers.StaffQ(r).FilterByUserID(userId).Get()
	if err != nil || resultStaff == nil {
		if resultStaff == nil {
			err = errors.New("resultStaff == nil")
		}
		return nil, err
	}
	resultPosition, err := helpers.PositionsQ(r).FilterByID(resultStaff.PositionID).Get()
	if err != nil || resultPosition == nil {
		if resultStaff == nil {
			err = errors.New("resultPosition == nil")
		}
		return nil, err
	}
	return resultPosition, nil
}

func GetPositionByUserHandler(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPositionByUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	resultPosition, err := GetPositionByUser(r, cast.ToInt64(request.UserKey.ID))

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
