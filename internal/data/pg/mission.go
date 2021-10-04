package pg

import (
	"database/sql"

	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableMission = "mission"

func NewMissionQ(db *pgdb.DB) data.MissionQ {
	return &missionQ{
		db:  db.Clone(),
		sql: squirrel.Select("*").From(tableMission),
	}
}

type missionQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func (d *missionQ) New() data.MissionQ {
	return NewMissionQ(d.db)
}

func (d *missionQ) Get() (*data.Mission, error) {
	var result data.Mission

	err := d.db.Get(&result, d.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (d *missionQ) Select(query pgdb.OffsetPageParams) ([]data.Mission, error) {
	var result []data.Mission

	err := d.db.Select(&result, query.ApplyTo( d.sql,"mission_id"))
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (d *missionQ) Insert(mission data.Mission) (data.Mission, error) {
	clauses := structs.Map(mission)

	query := squirrel.Insert(tableMission).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&mission, query)
	if err != nil {
		return data.Mission{}, err
	}

	return mission, err

}

func (d *missionQ) Update(mission data.Mission) (data.Mission, error) {
	clauses := structs.Map(mission)

	query := squirrel.Update(tableMission).Where(squirrel.Eq{"mission_id": mission.MissionId}).SetMap(clauses).Suffix("returning *")

	err := d.db.Get(&mission, query)
	if err != nil {
		return data.Mission{}, err
	}

	return mission, err
}

func (d missionQ) FilterById(missionId int64) data.MissionQ {
	d.sql = d.sql.Where(squirrel.Eq{"mission_id": missionId})

	return &d
}
