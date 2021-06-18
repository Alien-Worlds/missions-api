package config

import (
	"github.com/ethereum/go-ethereum/core/types"
	bscClient "github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	copTypes "gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/connectors/signed"
)

type Config interface {
	Block() Block
	BSC() bscClient.Client
	Contract() ContractAddress
	EventsConfig() EventsConfig

	pgdb.Databaser
	comfig.Logger
	copTypes.Copuser
	comfig.Listenerer
	signed.Clienter
}

type config struct {
	bsc bscClient.Client
	contract ContractAddress
	eventsConfig EventsConfig
	lastBlock types.Block

	pgdb.Databaser
	comfig.Logger
	copTypes.Copuser
	comfig.Listenerer
	signed.Clienter

	getter kv.Getter
	once   comfig.Once
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Clienter:   signed.NewClienter(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:  pgdb.NewDatabaser(getter),
	}
}
