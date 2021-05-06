package config

import (
	"github.com/binance-chain/bsc/core/types"
	bscClient "github.com/binance-chain/bsc/ethclient"
	"github.com/redcuckoo/bsc-checker-events/internal/config/hooks"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	BSC() bscClient.Client

	EventsConfig() EventsConfig

	Assets() []hooks.Assets

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

	assets []hooks.Assets

	lastBlock types.Block
}

func New(getter kv.Getter) Config {
	return &config{
		getter:    getter,
		Logger:    comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser: pgdb.NewDatabaser(getter),
	}
}
