package data

import (
	"time"

	"github.com/Digital-Voting-Team/staff-service/resources"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type StaffQ interface {
	New() StaffQ

	Get() (*Staff, error)
	Select() ([]Staff, error)

	Transaction(fn func(q StaffQ) error) error

	Insert(Staff) (Staff, error)
	Update(Staff) (Staff, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) StaffQ

	FilterByID(ids ...int64) StaffQ
	FilterByWorkStart(time time.Time) StaffQ
	FilterByWorkEnd(time time.Time) StaffQ
	FilterBySalaryUp(salaries ...float32) StaffQ
	FilterBySalaryBottom(salaries ...float32) StaffQ
	FilterByPosition(ids ...int64) StaffQ
	FilterByCafe(ids ...int64) StaffQ

	JoinPerson() StaffQ
	JoinPosition() StaffQ
}

type Staff struct {
	ID             int64                  `db:"id" structs:"-"`
	EmploymentDate *time.Time             `db:"employment_date" structs:"employment_date"`
	Salary         float32                `db:"salary" structs:"salary"`
	Status         resources.WorkerStatus `db:"status" structs:"status"`
	PersonID       int64                  `db:"person_id" structs:"person_id"`
	CafeID         int64                  `db:"cafe_id" structs:"cafe_id"`
	PositionID     int64                  `db:"position_id" structs:"position_id"`
	UserId         int64                  `db:"user_id" structs:"user_id"`
}
