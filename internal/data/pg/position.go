package pg

import (
	"database/sql"
	"gitlab.com/distributed_lab/kit/pgdb"
	"staff-service/internal/data"
	"staff-service/resources"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
)

const positionTableName = "public.position"

func NewPositionsQ(db *pgdb.DB) data.PositionsQ {
	return &positionQ{
		db:        db.Clone(),
		sql:       sq.Select("position.*").From(positionTableName),
		sqlUpdate: sq.Update(positionTableName).Suffix("returning *"),
	}
}

type positionQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *positionQ) New() data.PositionsQ {
	return NewPositionsQ(q.db)
}

func (q *positionQ) Get() (*data.Position, error) {
	var result data.Position
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *positionQ) Select() ([]data.Position, error) {
	var result []data.Position
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *positionQ) Update(position data.Position) (data.Position, error) {
	var result data.Position
	clauses := structs.Map(position)
	clauses["name"] = position.Name
	clauses["access_level"] = position.AccessLevel

	err := q.db.Get(&result, q.sqlUpdate.SetMap(clauses))

	return result, err
}

func (q *positionQ) Transaction(fn func(q data.PositionsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *positionQ) Insert(position data.Position) (data.Position, error) {
	clauses := structs.Map(position)
	clauses["name"] = position.Name
	clauses["access_level"] = position.AccessLevel

	var result data.Position
	stmt := sq.Insert(positionTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *positionQ) Delete(id int64) error {
	stmt := sq.Delete(positionTableName).Where(sq.Eq{"id": id})
	err := q.db.Exec(stmt)
	return err
}

func (q *positionQ) Page(pageParams pgdb.OffsetPageParams) data.PositionsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *positionQ) FilterByID(ids ...int64) data.PositionsQ {
	q.sql = q.sql.Where(sq.Eq{"id": ids})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"id": ids})
	return q
}

func (q *positionQ) FilterByName(names ...string) data.PositionsQ {
	q.sql = q.sql.Where(sq.Eq{"name": names})
	return q
}

func (q *positionQ) FilterByAccessLevel(accessLevels ...resources.AccessLevel) data.PositionsQ {
	q.sql = q.sql.Where(sq.Eq{"access_level": accessLevels})
	return q
}
