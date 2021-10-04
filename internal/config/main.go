package config

import (
	"net"

	bscClient "github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/connectors/signed"
)

type Config interface {
	Block() Block
	BSC() bscClient.Client
	Contract() ContractAddress
	EventsConfig() EventsConfig
	Listener() net.Listener

	pgdb.Databaser
	comfig.Logger
	// comfig.Listenerer
	signed.Clienter
}

type config struct {
	bsc bscClient.Client
	// contract ContractAddress
	eventsConfig EventsConfig
	// lastBlock types.Block
	// Listener() net.Listener
	pgdb.Databaser
	comfig.Logger
	// comfig.Listenerer
	signed.Clienter

	// getter kv.Getter
	// once   comfig.Once
}

func New(getter kv.Getter) Config {
	return &config{
		// getter:     getter,
		// Listenerer: comfig.NewListenerer(getter),
		Clienter:   signed.NewClienter(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:  pgdb.NewDatabaser(getter),
	}
}
