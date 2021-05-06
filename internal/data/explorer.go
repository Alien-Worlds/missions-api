package data

type ExplorerQ interface {
	New() ExplorerQ

	Get() (*Explorer, error)
	Select() ([]Explorer, error)

	Insert(mission Explorer) (Explorer, error)
	Update(mission Explorer) (Explorer, error)
}

type Explorer struct{
	ExplorerId    uint64
	totalStakeTLM uint64
	totalStakeBNB uint64
	missions []Mission
}