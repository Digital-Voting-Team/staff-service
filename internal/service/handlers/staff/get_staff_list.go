package handlers

import (
	"github.com/Digital-Voting-Team/staff-service/internal/data"
	"github.com/Digital-Voting-Team/staff-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/staff-service/internal/service/requests/staff"
	"github.com/Digital-Voting-Team/staff-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetStaffList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetStaffListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	staffQ := helpers.StaffQ(r)
	applyFilters(staffQ, request)
	staff, err := staffQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get staff")
		ape.Render(w, problems.InternalError())
		return
	}
	persons, err := helpers.PersonsQ(r).FilterByID(getPersonsIDs(staff)...).Select()
	positions, err := helpers.PositionsQ(r).FilterByID(getPositionsIDs(staff)...).Select()

	response := resources.StaffListResponse{
		Data:     newStaffList(staff),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newStaffIncluded(persons, positions),
	}
	ape.Render(w, response)
}

func applyFilters(q data.StaffQ, request requests.GetStaffListRequest) {
	q.Page(request.OffsetPageParams)

	if request.FilterWorkStart != nil {
		q.FilterByWorkStart(*request.FilterWorkStart)
	}

	if request.FilterWorkEnd != nil {
		q.FilterByWorkEnd(*request.FilterWorkEnd)
	}

	if len(request.FilterSalaryUp) > 0 {
		q.FilterBySalaryUp(request.FilterSalaryUp...)
	}

	if len(request.FilterSalaryBottom) > 0 {
		q.FilterBySalaryBottom(request.FilterSalaryBottom...)
	}

	if len(request.FilterPosition) > 0 {
		q.FilterByPosition(request.FilterPosition...)
	}

	if len(request.FilterCafe) > 0 {
		q.FilterByCafe(request.FilterCafe...)
	}
}

func newStaffList(staff []data.Staff) []resources.Staff {
	result := make([]resources.Staff, len(staff))
	for i, staff_ := range staff {
		result[i] = resources.Staff{
			Key: resources.NewKeyInt64(staff_.ID, resources.STAFF),
			Attributes: resources.StaffAttributes{
				EmploymentDate: *staff_.EmploymentDate,
				Salary:         staff_.Salary,
				Status:         &staff_.Status,
			},
			Relationships: resources.StaffRelationships{
				Person: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(staff_.PersonID, 10),
						Type: resources.PERSON,
					},
				},
				Position: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(staff_.PositionID, 10),
						Type: resources.POSITION,
					},
				},
				Cafe: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(staff_.CafeID, 10),
						Type: resources.CAFE_REF,
					},
				},
			},
		}
	}
	return result
}

func getPersonsIDs(staff []data.Staff) []int64 {
	personIDs := make([]int64, len(staff))
	for i := 0; i < len(staff); i++ {
		personIDs[i] = staff[i].PersonID
	}
	return personIDs
}

func getPositionsIDs(staff []data.Staff) []int64 {
	positionIds := make([]int64, len(staff))
	for i := 0; i < len(staff); i++ {
		positionIds[i] = staff[i].PositionID
	}
	return positionIds
}

func newStaffIncluded(persons []data.Person, positions []data.Position) resources.Included {
	result := resources.Included{}
	for _, item := range persons {
		resource := newPersonModel(item)
		result.Add(&resource)
	}
	for _, item := range positions {
		resource := newPositionModel(item)
		result.Add(&resource)
	}
	return result
}

func newPersonModel(person data.Person) resources.Person {
	return resources.Person{
		Key: resources.NewKeyInt64(person.ID, resources.PERSON),
		Attributes: resources.PersonAttributes{
			Name:     person.Name,
			Phone:    person.Phone,
			Email:    person.Email,
			Birthday: person.Birthday,
		},
		Relationships: resources.PersonRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(person.AddressID, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	}
}

func newPositionModel(position data.Position) resources.Position {
	return resources.Position{
		Key: resources.NewKeyInt64(position.ID, resources.POSITION),
		Attributes: resources.PositionAttributes{
			Name:        position.Name,
			AccessLevel: *position.AccessLevel,
		},
	}
}
