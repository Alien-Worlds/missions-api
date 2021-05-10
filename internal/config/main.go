package config

import (
	"github.com/ethereum/go-ethereum/core/types"
	//"github.com/binance-chain/bsc-static/bsc/common"
	//"github.com/binance-chain/bsc-static/bsc/core/types"
	//bscClient "github.com/binance-chain/bsc-static/bsc/ethclient"
	bscClient "github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	BSC() bscClient.Client

	EventsConfig() EventsConfig

	Contract() ContractAddress
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

	contract ContractAddress

	lastBlock types.Block
}

func New(getter kv.Getter) Config {
	return &config{
		getter:    getter,
		Logger:    comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser: pgdb.NewDatabaser(getter),
	}
}
