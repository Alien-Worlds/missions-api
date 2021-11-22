package data

type JoinEventQ interface {
	New() JoinEventQ

	Get() (*JoinEvent, error)
	Select() ([]JoinEvent, error)

	Insert(JoinEvent JoinEvent) (JoinEvent, error)
	Update(JoinEvent JoinEvent) (JoinEvent, error)

	FilterById(join_event_id string) JoinEventQ

}

type JoinEvent struct{
	TransactionId       string `db:"transaction_id" structs:"transaction_id"`
	// Explorer          string `db:"explorer_address" structs:"explorer_address"`
	// Mission           int64 `db:"mission" structs:"mission"`
	// // Withdrawn         bool `db:"withdrawn" structs:"withdrawn"`
	// NumberShips       int64 `db:"number_ships" structs:"number_ships"`
	// StakeTLM          int64 `db:"stake_tlm" structs:"stake_tlm"`
	// StakeBNB          uint64 `db:"stake_bnb" structs:"stake_bnb"`
}