package checker

import (
	//"binance-chain/bsc-static/bsc/common"
	"github.com/binance-chain/bsc-static/bsc/common"
	bscClient "github.com/binance-chain/bsc-static/bsc/ethclient"
	"github.com/redcuckoo/bsc-checker-events/internal/config"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"github.com/redcuckoo/bsc-checker-events/internal/data/pg"
	"gitlab.com/distributed_lab/logan/v3"
)

type Service struct {
	bscClient       bscClient.Client
	eventsConfig    config.EventsConfig
	missionQ        data.MissionQ
	explorerQ       data.ExplorerQ
	contractAddress common.Address
	lastBlockNumber uint64
	log             *logan.Entry
}

func New(cfg config.Config) *Service {
	log := cfg.Log().WithField("main_service", "checker")

	return &Service{
		bscClient: cfg.BSC(),
		eventsConfig: cfg.EventsConfig(),
		missionQ: pg.NewMissionQ(cfg.DB()),
		explorerQ: pg.NewExplorerQ(cfg.DB()),
		contractAddress: cfg.Contract(),
		lastBlockNumber: 0,
		log:    log,
	}
}
