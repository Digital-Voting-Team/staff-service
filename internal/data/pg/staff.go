package pg

import (
	"database/sql"
	"fmt"
	"github.com/Digital-Voting-Team/staff-service/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

const staffTableName = "public.staff"

func NewStaffQ(db *pgdb.DB) data.StaffQ {
	return &staffQ{
		db:        db.Clone(),
		sql:       sq.Select("staff.*").From(staffTableName),
		sqlUpdate: sq.Update(staffTableName).Suffix("returning *"),
	}
}

type staffQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (s *staffQ) New() data.StaffQ {
	return NewStaffQ(s.db)
}

func (s *staffQ) Get() (*data.Staff, error) {
	var result data.Staff
	err := s.db.Get(&result, s.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (s *staffQ) Select() ([]data.Staff, error) {
	var result []data.Staff
	err := s.db.Select(&result, s.sql)
	return result, err
}

func (s *staffQ) Transaction(fn func(q data.StaffQ) error) error {
	return s.db.Transaction(func() error {
		return fn(s)
	})
}

func (s *staffQ) Insert(staff data.Staff) (data.Staff, error) {
	clauses := structs.Map(staff)
	clauses["employment_date"] = staff.EmploymentDate
	clauses["salary"] = staff.Salary
	clauses["status"] = staff.Status
	clauses["person_id"] = staff.PersonID
	clauses["cafe_id"] = staff.CafeID
	clauses["position_id"] = staff.PositionID
	clauses["user_id"] = staff.PositionID

	var result data.Staff
	stmt := sq.Insert(staffTableName).SetMap(clauses).Suffix("returning *")
	err := s.db.Get(&result, stmt)

	return result, err
}

func (s *staffQ) Update(staff data.Staff) (data.Staff, error) {
	var result data.Staff
	clauses := structs.Map(staff)
	clauses["employment_date"] = staff.EmploymentDate
	clauses["salary"] = staff.Salary
	clauses["status"] = staff.Status
	clauses["person_id"] = staff.PersonID
	clauses["cafe_id"] = staff.CafeID
	clauses["position_id"] = staff.PositionID
	clauses["user_id"] = staff.PositionID

	err := s.db.Get(&result, s.sqlUpdate.SetMap(clauses))
	return result, err
}

func (s *staffQ) Delete(id int64) error {
	stmt := sq.Delete(staffTableName).Where(sq.Eq{"id": id})
	err := s.db.Exec(stmt)
	return err
}

func (s *staffQ) Page(pageParams pgdb.OffsetPageParams) data.StaffQ {
	s.sql = pageParams.ApplyTo(s.sql, "id")
	return s
}

func (s *staffQ) FilterByID(ids ...int64) data.StaffQ {
	s.sql = s.sql.Where(sq.Eq{"id": ids})
	s.sqlUpdate = s.sqlUpdate.Where(sq.Eq{"id": ids})
	return s
}

func (s *staffQ) FilterByWorkStart(time time.Time) data.StaffQ {
	stmt := sq.GtOrEq{"staff.employment_date": time}
	s.sql = s.sql.Where(stmt)
	// Will not work for update
	// s.sqlUpdate = s.sqlUpdate.Where(stmt)
	return s
}

func (s *staffQ) FilterByWorkEnd(time time.Time) data.StaffQ {
	stmt := sq.LtOrEq{"staff.employment_date": time}
	s.sql = s.sql.Where(stmt)
	// Will not work for update
	// s.sqlUpdate = s.sqlUpdate.Where(stmt)
	return s
}

func (s *staffQ) FilterBySalaryUp(salaries ...float32) data.StaffQ {
	stmt := sq.LtOrEq{"staff.salary": salaries}
	s.sql = s.sql.Where(stmt)
	return s
}

func (s *staffQ) FilterBySalaryBottom(salaries ...float32) data.StaffQ {
	stmt := sq.GtOrEq{"staff.salary": salaries}
	s.sql = s.sql.Where(stmt)
	return s
}

func (s *staffQ) FilterByPosition(ids ...int64) data.StaffQ {
	s.sql = s.sql.Where(sq.Eq{"staff.position_id": ids})
	return s
}

func (s *staffQ) FilterByCafe(ids ...int64) data.StaffQ {
	s.sql = s.sql.Where(sq.Eq{"staff.cafe_id": ids})
	return s
}

func (s *staffQ) JoinPosition() data.StaffQ {
	stmt := fmt.Sprintf("%s as staff on public.position.id = staff.position_id",
		staffTableName)
	s.sql = s.sql.Join(stmt)
	return s
}

func (s *staffQ) JoinPerson() data.StaffQ {
	stmt := fmt.Sprintf("%s as staff on public.person.id = staff.person_id",
		staffTableName)
	s.sql = s.sql.Join(stmt)
	return s
}
