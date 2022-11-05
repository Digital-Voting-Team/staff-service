package handlers

import (
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/staff"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetStaff(w http.ResponseWriter, r *http.Request) {
	accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	if *accessLevel < resources.Manager {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewGetStaffRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resultStaff, err := helpers.StaffQ(r).FilterByID(request.StaffID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get staff from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if resultStaff == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	relatePerson, err := helpers.PersonsQ(r).FilterByID(resultStaff.PersonID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get person")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	relatePosition, err := helpers.PositionsQ(r).FilterByID(resultStaff.PositionID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get position")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var includes resources.Included
	includes.Add(&resources.Person{
		Key: resources.NewKeyInt64(relatePerson.ID, resources.PERSON),
		Attributes: resources.PersonAttributes{
			Name:     relatePerson.Name,
			Phone:    relatePerson.Phone,
			Email:    relatePerson.Email,
			Birthday: relatePerson.Birthday,
		},
		Relationships: resources.PersonRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relatePerson.AddressID, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	})

	includes.Add(&resources.Position{
		Key: resources.NewKeyInt64(relatePosition.ID, resources.POSITION),
		Attributes: resources.PositionAttributes{
			Name:        relatePosition.Name,
			AccessLevel: *relatePosition.AccessLevel,
		},
	})

	result := resources.StaffResponse{
		Data: resources.Staff{
			Key: resources.NewKeyInt64(resultStaff.ID, resources.STAFF),
			Attributes: resources.StaffAttributes{
				EmploymentDate: *resultStaff.EmploymentDate,
				Salary:         resultStaff.Salary,
				Status:         &resultStaff.Status,
			},
			Relationships: resources.StaffRelationships{
				Person: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultStaff.PersonID, 10),
						Type: resources.PERSON,
					},
				},
				Position: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultStaff.PositionID, 10),
						Type: resources.POSITION,
					},
				},
				Cafe: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultStaff.CafeID, 10),
						Type: resources.CAFE_REF,
					},
				},
				User: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultStaff.UserId, 10),
						Type: resources.USER_REF,
					},
				},
			},
		},
		Included: includes,
	}
	ape.Render(w, result)
}
