package user

import (
	endpoints2 "github.com/Digital-Voting-Team/auth-serivce/endpoints"
	"github.com/Digital-Voting-Team/staff-service/internal/config"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetPositionByJWT(endpointsConf *config.EndpointsConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtResponse, err := endpoints2.ValidateToken(
			r.Header.Get("Authorization"),
			endpointsConf.Endpoints["auth-service"],
		)
		if err != nil || jwtResponse.Data.ID == "" {
			helpers.Log(r).WithError(err).Info("auth failed")
			ape.Render(w, problems.BadRequest(err))
			return
		}

		resultPosition, err := GetPositionByUser(r, cast.ToInt64(jwtResponse.Data.Relationships.User.Data.ID))
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
}
