package handlers

import (
	"github.com/Digital-Voting-Team/staff-service/internal/data"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/staff"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"github.com/spf13/cast"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateStaff(w http.ResponseWriter, r *http.Request) {
	accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	if *accessLevel < resources.Manager {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewCreateStaffRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	staff := data.Staff{
		EmploymentDate: &request.Data.Attributes.EmploymentDate,
		Salary:         request.Data.Attributes.Salary,
		Status:         *request.Data.Attributes.Status,
		PersonID:       cast.ToInt64(request.Data.Relationships.Person.Data.ID),
		CafeID:         cast.ToInt64(request.Data.Relationships.Cafe.Data.ID),
		PositionID:     cast.ToInt64(request.Data.Relationships.Position.Data.ID),
		UserId:         cast.ToInt64(request.Data.Relationships.User.Data.ID),
	}

	relatePerson, err := helpers.PersonsQ(r).FilterByID(staff.PersonID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get person")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	relatePosition, err := helpers.PositionsQ(r).FilterByID(staff.PositionID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get position")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resultStaffByUser, err := helpers.StaffQ(r).FilterByUserID(staff.UserId).Get()
	if resultStaffByUser != nil {
		helpers.Log(r).WithError(err).Error("user already related to staff")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	var resultStaff data.Staff
	resultStaff, err = helpers.StaffQ(r).Insert(staff)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create staff")
		ape.RenderErr(w, problems.InternalError())
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
