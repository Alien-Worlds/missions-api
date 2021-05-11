package data

type MissionQ interface {
	New() MissionQ

	Get() (*Mission, error)
	Select() ([]Mission, error)

	Insert(mission Mission) (Mission, error)
	Update(mission Mission) (Mission, error)

	FilterById(missionId uint64) MissionQ
}

type Mission struct{
	MissionId     uint64 `db:"mission_id" structs:"mission_id"`
	Description   string `db:"description" structs:"description"`
	Name          string `db:"name" structs:"name"`
	BoardingTime  int64 `db:"boarding_time" structs:"boarding_time"`
	LaunchTime    int64 `db:"launch_time" structs:"launch_time"`
	EndTime       int64 `db:"end_time" structs:"end_time"`
	Duration      int64 `db:"duration" structs:"duration"`
	MissionType   int64 `db:"mission_type" structs:"mission_type"`
	Reward        int64 `db:"reward" structs:"reward"`
	SpaceshipCost int64 `db:"spaceship_cost" structs:"spaceship_cost"`
	MissionPower  int64 `db:"mission_power" structs:"mission_power"`
	TotalShips    int64 `db:"total_ships" structs:"total_ships"`
	NftContract   []byte `db:"nft_contract" structs:"nft_contract"`
	NftTokenURI   string `db:"nft_token_uri" structs:"nft_token_uri"`
}