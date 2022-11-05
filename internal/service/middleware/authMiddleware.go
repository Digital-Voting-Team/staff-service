package middleware

import (
	"context"
	"github.com/Digital-Voting-Team/auth-serivce/endpoints"
	"github.com/Digital-Voting-Team/staff-service/internal/config"
	"github.com/Digital-Voting-Team/staff-service/internal/service/handlers/user"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func BasicAuth(endpointsConf *config.EndpointsConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtResponse, err := endpoints.ValidateToken(
				r.Header.Get("Authorization"),
				endpointsConf.Endpoints["auth-service"],
			)
			if jwtResponse == nil {
				helpers.Log(r).WithError(err).Info("auth failed, jwtResponse == nil")
				ape.Render(w, problems.BadRequest(err))
				return
			}
			if err != nil || jwtResponse.Data.ID == "" {
				helpers.Log(r).WithError(err).Info("auth failed")
				ape.Render(w, problems.BadRequest(err))
				return
			}
			position, err := user.GetPositionByUser(r, cast.ToInt64(jwtResponse.Data.Relationships.User.Data.ID))
			if position == nil {
				helpers.Log(r).WithError(err).Info("auth failed, position == nil")
				ape.Render(w, problems.BadRequest(err))
				return
			}
			if err != nil || jwtResponse.Data.ID == "" {
				helpers.Log(r).WithError(err).Info("auth failed, no staff to user")
				ape.Render(w, problems.Forbidden())
				return
			}
			ctx := context.WithValue(r.Context(), "accessLevel", position.AccessLevel)
			ctx = context.WithValue(ctx, "userId", cast.ToInt64(jwtResponse.Data.Relationships.User.Data.ID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
