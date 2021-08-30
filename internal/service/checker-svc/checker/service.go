package checker

import (
	"context"
	"github.com/Alien-Worlds/missions-api/internal/contracts"
	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"math/big"
	"time"
)

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Running checker service")

	running.WithBackOff(ctx, s.log, "new-checker-service", func(ctx context.Context) error {
		if s.lastBlockNumber == 0 {
			currentBlock, err := s.bscClient.BlockByNumber(ctx, nil)

			if err != nil {
				return errors.Wrap(err, "failed to fetch current block")
			}

			s.lastBlockNumber = currentBlock.NumberU64()
			s.log.Info("Fetched start block successfully")
			s.log.Infof("Current block ", currentBlock.Number())
			return nil
		}

		currBlock, err := s.bscClient.BlockByNumber(ctx, nil)

		if err != nil {
			return errors.Wrap(err, "failed to fetch current block")
		}

		if s.lastBlockNumber > currBlock.NumberU64() {
			return nil
		}

		//the limit for filtering logs is 500 blocks
		var toBlockNum uint64

		if currBlock.NumberU64() - s.lastBlockNumber > 500{
			toBlockNum = s.lastBlockNumber + 500

			s.log.Infof("Fetching events from block %d to block %d",s.lastBlockNumber, toBlockNum)
		}else{
			toBlockNum = currBlock.NumberU64()

			s.log.Infof("Fetching events from block %d to head of blockchain,  block %d",s.lastBlockNumber, toBlockNum)
		}

		logs, err := s.bscClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(s.lastBlockNumber)),
			ToBlock:  big.NewInt(int64(toBlockNum)),
			Addresses: []common.Address{
				common.HexToAddress(s.contractAddress.Address),
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to set filter query for filter logs")
		}

		if currBlock.NumberU64() - s.lastBlockNumber > 500{
			s.lastBlockNumber = s.lastBlockNumber + 501
		}else{
			s.lastBlockNumber = currBlock.NumberU64() + 1
		}

		err = s.process(ctx, logs)
		if err != nil {
			return errors.Wrap(err, "failed to process events")
		}

		return nil
	}, time.Second, 30*time.Second, time.Minute)
}

func (s *Service) process(ctx context.Context, logs []types.Log) error {
	s.log.Infof("Running process(), logs amount %d", len(logs))
	for _, event := range logs {
		err := s.processEvent(ctx, event)
		if err != nil {
			s.log.WithError(err).Error("failed to process event")
			continue
		}
	}

	return nil
}

func (s *Service) processEvent(ctx context.Context, event types.Log) error {
	spaceshipToken, err := SpaceshipStaking.NewSpaceshipStaking(common.HexToAddress(s.contractAddress.Address), &s.bscClient)
	if err != nil {
		return errors.Wrap(err, "failed init new SpaceshipStaking")
	}

	if event.Topics == nil {
		return nil
	}

	if event.Topics[0].String() != s.eventsConfig.MissionCreatedHash &&
		event.Topics[0].String() != s.eventsConfig.MissionJoinedHash &&
		event.Topics[0].String() != s.eventsConfig.RewardWithdrawnHash {
		return nil
	}

	switch event.Topics[0].String() {
	case s.eventsConfig.MissionCreatedHash:
		err = s.processMissionCreated(event, spaceshipToken)

		if err != nil {
			return errors.Wrap(err, "failed process mission created event")
		}
	case s.eventsConfig.MissionJoinedHash:
		err = s.processMissionJoined(event, spaceshipToken)

		if err != nil {
			return errors.Wrap(err, "failed process mission joined event")
		}
	case s.eventsConfig.RewardWithdrawnHash:
		err = s.processRewardWithdrawn(event, spaceshipToken)

		if err != nil {
			return errors.Wrap(err, "failed process reward withdrawn event")
		}
	default:
		return nil
	}

	return nil
}

func (s *Service) processMissionCreated(event types.Log, spaceshipToken *SpaceshipStaking.SpaceshipStaking) error {
	s.log.Info("Started parsing MissionCreated event")

	missionCreated, err := spaceshipToken.ParseMissionCreated(event)

	if err != nil {
		return errors.Wrap(err, "failed parse mission created")
	}

	mission, err := spaceshipToken.Missions(&bind.CallOpts{}, missionCreated.Id)

	if err != nil {
		return errors.Wrap(err, "failed get info about mission created")
	}

	missionDB := data.Mission{
		MissionId:     missionCreated.Id.Uint64(),
		Description:   mission.Description,
		Name:          mission.Name,
		BoardingTime:  int64(mission.BoardingTime),
		LaunchTime:    int64(mission.LaunchTime),
		EndTime:       int64(mission.Duration + mission.LaunchTime),
		Duration:      int64(mission.Duration),
		MissionType:   int64(mission.MissionType),
		Reward:        mission.Reward.Int64(),
		SpaceshipCost: mission.SpaceshipCost.Int64(),
		MissionPower:  mission.MissionPower.Int64(),
		TotalShips:    0,
		NftContract:   mission.NftInfo.ContractAddress.Bytes(),
		NftTokenURI:   mission.NftInfo.TokenURI,
	}

	missionCheck, err := s.missionQ.FilterById(int64(missionDB.MissionId)).Get()

	if err != nil {
		return errors.Wrap(err, "failed to get from db, mission")
	}

	if missionCheck == nil {
		_, err = s.missionQ.Insert(missionDB)

		if err != nil {
			return errors.Wrap(err, "failed to insert to db, mission")
		}
	} else {
		_, err = s.missionQ.Update(missionDB)

		if err != nil {
			return errors.Wrap(err, "failed to update db, mission")
		}
	}

	s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission created event")
	return nil
}

func (s *Service) processMissionJoined(event types.Log, spaceshipToken *SpaceshipStaking.SpaceshipStaking) error {
	s.log.Info("Started parsing MissionJoined event")
	missionJoined, err := spaceshipToken.ParseMissionJoined(event)
	if err != nil {
		return errors.Wrap(err, "failed parse mission joined")
	}

	investINFO, err := spaceshipToken.MissionToUsersInvest(&bind.CallOpts{}, missionJoined.MissionId, missionJoined.Player)

	if err != nil {
		return errors.Wrap(err, "failed get info about mission joined")
	}

	missionFromContract, err := spaceshipToken.Missions(&bind.CallOpts{}, missionJoined.MissionId)

	if err != nil {
		return errors.Wrap(err, "failed fetch info from blockchain, mission")
	}

	explorer, err := s.explorerQ.FilterByAddress(missionJoined.Player.String()).Get()

	if err != nil {
		return errors.Wrap(err, "failed fetch info from db, table explorer, explorer")
	}

	//updating explorer-mission table
	var explorerMissionDB *data.ExplorerMission
	explorerMissionDB = nil

	if explorer != nil {
		explorerMissionDB, err = s.explorerMissionQ.FilterByMission(missionJoined.MissionId.Int64()).FilterByExplorer(int64(explorer.ExplorerId)).Get()
	}

	if err != nil {
		return errors.Wrap(err, "failed update db, explorer-mission")
	}

	//updating mission table
	missionDB, err := s.missionQ.FilterById(missionJoined.MissionId.Int64()).Get()

	if err != nil {
		return errors.Wrap(err, "failed fetch info from db, table mission, mission")
	}

	if missionDB != nil {
		missionDB.MissionPower = missionFromContract.MissionPower.Int64()

		if explorerMissionDB == nil {
			missionDB.TotalShips += investINFO.Ships.Int64()
		} else {
			missionDB.TotalShips = (missionDB.TotalShips - explorerMissionDB.NumberShips) + investINFO.Ships.Int64()
		}

		_, err = s.missionQ.Update(*missionDB)

		if err != nil {
			return errors.Wrap(err, "failed update db, mission")
		}
	} else {
		return errors.Wrap(err, "failed get info from db about joining mission, resynchronize")
	}

	//update table explorer
	explorerDB := data.Explorer{
		ExplorerAddress: missionJoined.Player.String(),
		TotalStakeTLM:   missionFromContract.SpaceshipCost.Int64() * investINFO.Ships.Int64(),
		TotalStakeBNB:   investINFO.BNBAmount.Int64(),
	}

	if explorer != nil {
		if explorerMissionDB == nil {
			explorerDB.TotalStakeTLM = explorer.TotalStakeTLM + missionFromContract.SpaceshipCost.Int64()*investINFO.Ships.Int64()
			explorerDB.TotalStakeBNB = explorer.TotalStakeBNB + investINFO.BNBAmount.Int64()
		} else {
			explorerDB.TotalStakeTLM = (explorer.TotalStakeTLM - explorerMissionDB.TotalStakeTLM) + missionFromContract.SpaceshipCost.Int64()*investINFO.Ships.Int64()
			explorerDB.TotalStakeBNB = (explorer.TotalStakeBNB - explorerMissionDB.TotalStakeBNB) + investINFO.BNBAmount.Int64()
		}

		_, err = s.explorerQ.Update(explorerDB)

		if err != nil {
			return errors.Wrap(err, "failed update db, explorer")
		}
	} else {
		_, err = s.explorerQ.Insert(explorerDB)

		if err != nil {
			return errors.Wrap(err, "failed insert to db, explorer")
		}
	}

	explorerDBNew, err := s.explorerQ.FilterByAddress(explorerDB.ExplorerAddress).Get()

	if err != nil || explorerDBNew == nil {
		return errors.Wrap(err, "failed fetch inserted explorer from database, explorer")
	}

	if explorerMissionDB == nil {
		explorerMissionDB := data.ExplorerMission{
			Explorer:      int64(explorerDBNew.ExplorerId),
			Mission:       missionJoined.MissionId.Int64(),
			Withdrawn:     false,
			NumberShips:   investINFO.Ships.Int64(),
			TotalStakeTLM: investINFO.Ships.Int64() * missionFromContract.SpaceshipCost.Int64(),
			TotalStakeBNB: investINFO.BNBAmount.Int64(),
		}

		_, err = s.explorerMissionQ.Insert(explorerMissionDB)

		if err != nil {
			return errors.Wrap(err, "failed insert to db, explorer-mission")
		}
	} else {
		explorerMissionDB.NumberShips = investINFO.Ships.Int64()
		explorerMissionDB.TotalStakeTLM = explorerMissionDB.NumberShips * missionFromContract.SpaceshipCost.Int64()
		explorerMissionDB.TotalStakeBNB = investINFO.BNBAmount.Int64()

		_, err = s.explorerMissionQ.Update(*explorerMissionDB)

		if err != nil {
			return errors.Wrap(err, "failed update db, explorer-mission")
		}
	}

	s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission joined")

	return nil
}

func (s *Service) processRewardWithdrawn(event types.Log, spaceshipToken *SpaceshipStaking.SpaceshipStaking) error {
	s.log.Info("Started parsing RewardWithdrawn event")

	rewardWithdrawn, err := spaceshipToken.ParseRewardWithdrawn(event)

	if err != nil {
		return errors.Wrap(err, "failed parse reward withdrawn")
	}

	missionDB, err := s.missionQ.FilterById(rewardWithdrawn.MissionId.Int64()).Get()

	if err != nil {
		return errors.Wrap(err, "failed fetch info from db, table mission, mission")
	}

	if missionDB == nil {
		return errors.New("failed record withdraw, such mission is not in db, resynchronize")
	}

	explorerDB, err := s.explorerQ.FilterByAddress(rewardWithdrawn.Player.String()).Get()

	if err != nil {
		return errors.Wrap(err, "failed fetch info from db, table explorer, explorer")
	}

	if explorerDB == nil {
		return errors.New("failed record withdraw, such explorer is not in db, resynchronize")
	}

	explorerMissionDB, err := s.explorerMissionQ.FilterByMission(rewardWithdrawn.MissionId.Int64()).FilterByExplorer(int64(explorerDB.ExplorerId)).Get()

	if err != nil {
		return errors.Wrap(err, "failed get info about mission-explorer")
	}

	if explorerMissionDB == nil {
		return errors.New("failed record withdraw, such explorer-mission connection is not in db, resynchronize")
	} else {
		explorerMissionDB.Withdrawn = true

		_, err = s.explorerMissionQ.Update(*explorerMissionDB)

		if err != nil {
			return errors.Wrap(err, "failed to update db, explorer-mission")
		}
	}

	s.log.WithField("mission_id", explorerMissionDB.Mission).Info("Success get reward withdrawn")

	return nil
}
