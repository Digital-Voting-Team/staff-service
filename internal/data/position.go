package data

import (
	"github.com/Digital-Voting-Team/staff-service/resources"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type PositionsQ interface {
	New() PositionsQ

	Get() (*Position, error)
	Select() ([]Position, error)

	Transaction(fn func(q PositionsQ) error) error

	Insert(Position) (Position, error)
	Update(Position) (Position, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) PositionsQ

	FilterByID(ids ...int64) PositionsQ
	FilterByName(names ...string) PositionsQ
	FilterByAccessLevel(accessLevels ...resources.AccessLevel) PositionsQ
}

type Position struct {
	ID          int64                  `db:"id" structs:"-"`
	Name        string                 `db:"name" structs:"name"`
	AccessLevel *resources.AccessLevel `db:"access_level" structs:"access_level"`
}
