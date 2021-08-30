package checker

import (
	"github.com/Alien-Worlds/missions-api/internal/config"
	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/Alien-Worlds/missions-api/internal/data/pg"
	bscClient "github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
)

type Service struct {
	bscClient        bscClient.Client
	eventsConfig     config.EventsConfig
	missionQ         data.MissionQ
	explorerQ        data.ExplorerQ
	explorerMissionQ data.ExplorerMissionQ
	contractAddress  config.ContractAddress
	lastBlockNumber  uint64
	log              *logan.Entry
	db               *pgdb.DB
}

func New(cfg config.Config) *Service {
	log := cfg.Log().WithField("main_service", "checker")

	return &Service{
		bscClient:        cfg.BSC(),
		eventsConfig:     cfg.EventsConfig(),
		missionQ:         pg.NewMissionQ(cfg.DB()),
		explorerQ:        pg.NewExplorerQ(cfg.DB()),
		explorerMissionQ: pg.NewExplorerMissionQ(cfg.DB()),
		contractAddress:  cfg.Contract(),
		lastBlockNumber:  uint64(cfg.Block().FromBlockNum),
		log:              log,
		db:               cfg.DB(),
	}
}
