package data

type ExplorerMissionQ interface {
	New() ExplorerMissionQ

	Get() (*ExplorerMission, error)
	Select() ([]ExplorerMission, error)

	Insert(explorerMission ExplorerMission) (ExplorerMission, error)
	Update(explorerMission ExplorerMission) (ExplorerMission, error)
}

type ExplorerMission struct{
	ExplorerMissionId uint64
	Explorer          string
	Mission           uint64
	Withdrawn         bool
	NumberShips       uint64
	TotalStakeTLM     uint64
	TotalStakeBNB     uint64
}