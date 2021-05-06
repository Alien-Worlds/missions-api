package config

import (
	"github.com/binance-chain/bsc/common"
	"github.com/binance-chain/bsc/core/types"
	bscClient "github.com/binance-chain/bsc/ethclient"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	BSC() bscClient.Client

	EventsConfig() EventsConfig


	pgdb.Databaser
	comfig.Logger
}

type config struct {
	bsc bscClient.Client

	eventsConfig EventsConfig

	pgdb.Databaser
	comfig.Logger

	getter kv.Getter
	once   comfig.Once

	contract common.Address

	lastBlock types.Block
}

func New(getter kv.Getter) Config {
	return &config{
		getter:    getter,
		Logger:    comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser: pgdb.NewDatabaser(getter),
	}
}
