package data

type ExplorerQ interface {
	New() ExplorerQ

	Get() (*Explorer, error)

	Select() ([]Explorer, error)

	Insert(explorer Explorer) (Explorer, error)
	Update(explorer Explorer) (Explorer, error)
	
	FilterByAddress(explorerAddress string) ExplorerQ
}

type Explorer struct{
	ExplorerId    uint64 `db:"explorer_id" structs:"-"`
	ExplorerAddress string `db:"explorer_address" structs:"explorer_address"`
	TotalStakeTLM int64 `db:"total_stake_tlm" structs:"total_stake_tlm"`
	TotalStakeBNB int64 `db:"total_stake_bnb" structs:"total_stake_bnb"`
}