package data

type ExplorerQ interface {
	New() ExplorerQ

	Get() (*Explorer, error)
	SelectTotalTLM(explorerId string) (*Explorer, error)
	SelectTotalBNB(explorerId string) (*Explorer, error)

	Select() ([]Explorer, error)

	Insert(mission Explorer) (Explorer, error)
	Update(mission Explorer) (Explorer, error)
}

type Explorer struct{
	ExplorerId    string
	TotalStakeTLM uint64
	TotalStakeBNB uint64
}