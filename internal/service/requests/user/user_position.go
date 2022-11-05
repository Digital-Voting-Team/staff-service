package user

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

type GetPositionByUserRequest struct {
	UserKey resources.Key
}

func NewGetPositionByUserRequest(r *http.Request) (GetPositionByUserRequest, error) {
	var request GetPositionByUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
