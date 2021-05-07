package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableExplorerMission = "ExplorerMission"

func NewExplorerMissionQ(db *pgdb.DB) data.ExplorerMissionQ{
	return &explorerMissionQ{
		db:  db.Clone(),
		sql: squirrel.Select("*").From(tableExplorer),
	}
}

type explorerMissionQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func (d *explorerMissionQ) New() data.ExplorerMissionQ {
	return NewExplorerMissionQ(d.db)
}

func (d *explorerMissionQ) Get() (*data.ExplorerMission, error) {
	var result data.ExplorerMission
	err := d.db.Exec(d.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (d *explorerMissionQ) Select() ([]data.ExplorerMission, error) {
	var result []data.ExplorerMission
	err := d.db.Select(&result, d.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (d *explorerMissionQ) Insert(explorerMission data.ExplorerMission) (data.ExplorerMission, error) {
	clauses := structs.Map(explorerMission)

	query := squirrel.Insert(tableExplorerMission).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&explorerMission, query)
	if err != nil {
		return data.ExplorerMission{}, err
	}

	return explorerMission, err

}

func (d *explorerMissionQ) Update(explorerMission data.ExplorerMission) (data.ExplorerMission, error) {
	clauses := structs.Map(explorerMission)

	query := squirrel.Update(tableExplorerMission).Where(squirrel.Eq{"id": explorerMission.ExplorerMissionId}).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&explorerMission, query)
	if err != nil {
		return data.ExplorerMission{}, err
	}

	return explorerMission, err
}