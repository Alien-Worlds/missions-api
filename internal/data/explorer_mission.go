package data

type ExplorerMissionQ interface {
	New() ExplorerMissionQ

	Get() (*ExplorerMission, error)
	Select() ([]ExplorerMission, error)

	Insert(explorerMission ExplorerMission) (ExplorerMission, error)
	Update(explorerMission ExplorerMission) (ExplorerMission, error)

	FilterByMission(missionId uint64) ExplorerMissionQ
	FilterByExplorer(explorerId uint64) ExplorerMissionQ
}

type ExplorerMission struct{
	ExplorerMissionId uint64 `db:"explorer_mission_id" structs:"-"`
	Explorer          int64 `db:"explorer" structs:"explorer"`
	Mission           int64 `db:"mission" structs:"mission"`
	Withdrawn         bool `db:"withdrawn" structs:"withdrawn"`
	NumberShips       int64 `db:"number_ships" structs:"number_ships"`
	TotalStakeTLM     int64 `db:"total_stake_tlm" structs:"total_stake_tlm"`
	TotalStakeBNB     int64 `db:"total_stake_bnb" structs:"total_stake_bnb"`
}