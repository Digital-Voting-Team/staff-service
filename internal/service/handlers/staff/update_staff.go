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

func UpdateStaff(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateStaffRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	staff, err := helpers.StaffQ(r).FilterByID(request.StaffID).Get()
	if staff == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	userId := r.Context().Value("userId").(int64)
	accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	staffId, _, _, err := helpers.GetIdsForGivenUser(r, userId)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong relations")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if *accessLevel != resources.Admin && staffId != staff.ID {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	newStaff := data.Staff{
		EmploymentDate: &request.Data.Attributes.EmploymentDate,
		Salary:         request.Data.Attributes.Salary,
		Status:         *request.Data.Attributes.Status,
		PersonID:       cast.ToInt64(request.Data.Relationships.Person.Data.ID),
		CafeID:         cast.ToInt64(request.Data.Relationships.Cafe.Data.ID),
		PositionID:     cast.ToInt64(request.Data.Relationships.Position.Data.ID),
		UserId:         cast.ToInt64(request.Data.Relationships.User.Data.ID),
	}

	resultStaffByPerson, err := helpers.StaffQ(r).FilterByPersonID(staff.PersonID).Get()
	if resultStaffByPerson != nil {
		helpers.Log(r).WithError(err).Error("person already related to staff")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	relatePerson, err := helpers.PersonsQ(r).FilterByID(newStaff.PersonID).Get()
	if err != nil || relatePerson == nil {
		helpers.Log(r).WithError(err).Error("failed to get new person")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	relatePosition, err := helpers.PositionsQ(r).FilterByID(newStaff.PersonID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get position")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if staff.UserId != newStaff.UserId {
		helpers.Log(r).WithError(err).Error("cannot change staff to user relation")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	resultStaffByUser, err := helpers.StaffQ(r).FilterByUserID(staff.UserId).Get()
	if resultStaffByUser == nil {
		helpers.Log(r).WithError(err).Error("user not assigned to staff")
		ape.RenderErr(w, problems.Conflict())
		return
	}
	if resultStaffByUser.ID == 0 || resultStaffByUser.UserId != newStaff.UserId {
		helpers.Log(r).WithError(err).Error("invalid user to update")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	var resultStaff data.Staff
	resultStaff, err = helpers.StaffQ(r).FilterByID(staff.ID).Update(newStaff)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update staff")
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
