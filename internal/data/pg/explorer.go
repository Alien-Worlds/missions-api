package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableExplorer = "Explorer"

func NewExplorerQ(db *pgdb.DB) data.ExplorerQ{
	return &explorerQ{
		db:  db.Clone(),
		sql: squirrel.Select("*").From(tableExplorer),
	}
}

type explorerQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func (d *explorerQ) New() data.ExplorerQ {
	return NewExplorerQ(d.db)
}

func (d *explorerQ) Get() (*data.Explorer, error) {
	var result data.Explorer
	err := d.db.Exec(d.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (d *explorerQ) Select() ([]data.Explorer, error) {
	var result []data.Explorer
	err := d.db.Select(&result, d.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (d *explorerQ) Insert(explorer data.Explorer) (data.Explorer, error) {
	clauses := structs.Map(explorer)

	query := squirrel.Insert(tableExplorer).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&explorer, query)
	if err != nil {
		return data.Explorer{}, err
	}

	return explorer, err

}

func (d *explorerQ) Update(explorer data.Explorer) (data.Explorer, error) {
	clauses := structs.Map(explorer)

	query := squirrel.Update(tableExplorer).Where(squirrel.Eq{"id": explorer.ExplorerId}).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&explorer, query)
	if err != nil {
		return data.Explorer{}, err
	}

	return explorer, err
}
