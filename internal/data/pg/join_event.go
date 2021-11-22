package pg

import (
	"database/sql"

	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableJoinEvent = "join_event"

func NewJoinEventQ(db *pgdb.DB) data.JoinEventQ {
	return &joinEventQ{
		db:  db.Clone(),
		sql: squirrel.Select("*").From(tableJoinEvent).OrderBy("transaction_id DESC"),
	}
}

type joinEventQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func (d *joinEventQ) New() data.JoinEventQ {
	return NewJoinEventQ(d.db)
}

func (d *joinEventQ) Get() (*data.JoinEvent, error) {
	var result data.JoinEvent
	err := d.db.Get(&result, d.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (d *joinEventQ) Select() ([]data.JoinEvent, error) {
	var result []data.JoinEvent
	err := d.db.Select(&result, d.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (d *joinEventQ) Insert(joinEvent data.JoinEvent) (data.JoinEvent, error) {
	clauses := structs.Map(joinEvent)

	query := squirrel.Insert(tableJoinEvent).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&joinEvent, query)
	if err != nil {
		return data.JoinEvent{}, err
	}

	return joinEvent, err

}

func (d *joinEventQ) Update(joinEvent data.JoinEvent) (data.JoinEvent, error) {
	clauses := structs.Map(joinEvent)

	query := squirrel.Update(tableJoinEvent).Where(squirrel.Eq{"transaction_id": joinEvent.TransactionId}).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&joinEvent, query)
	if err != nil {
		return data.JoinEvent{}, err
	}

	return joinEvent, err
}

func (d joinEventQ) FilterById(transaction_id string) data.JoinEventQ {
	d.sql = d.sql.Where(squirrel.Like{"transaction_id": transaction_id})
	return &d
}

