package data

type MissionQ interface {
	New() MissionQ

	Get() (*Mission, error)
	Select() ([]Mission, error)

	Insert(mission Mission) (Mission, error)
	Update(mission Mission) (Mission, error)
}

type Mission struct{
	MissionId     uint64
	Description   string
	Name          string
	BoardingTime  uint64
	LaunchTime    uint64
	EndTime       uint64
	Duration      uint64
	MissionType   uint64
	Reward        uint64
	SpaceshipCost uint64
	MissionPower  uint64
	TotalShips    uint64
	NftContract   []byte
	NftTokenURI   string
}