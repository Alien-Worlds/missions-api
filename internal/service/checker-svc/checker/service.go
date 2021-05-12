package checker

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	//"math/big"
	"time"

	//_ "github.com/binance-chain/bsc-static/bsc"
	//"github.com/binance-chain/bsc-static/bsc/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//_ "github.com/binance-chain/bsc-static/bsc/common"
	"github.com/ethereum/go-ethereum/core/types"
	//"github.com/binance-chain/bsc-static/bsc/core/types"
	//_ "github.com/lib/pq"
	"github.com/redcuckoo/bsc-checker-events/internal/contracts"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Running checker service")

	running.WithBackOff(ctx, s.log, "new-checker-service", func(ctx context.Context) error {
		//if s.lastBlockNumber == 0 {
		//	currentBlock, err := s.bscClient.BlockByNumber(ctx, nil)
		//
		//	if err != nil {
		//		return errors.Wrap(err, "failed to fetch current block")
		//	}
		//
		//	s.lastBlockNumber = currentBlock.NumberU64()
		//	s.log.Info("Fetched start block successfully")
		//	s.log.Infof("Current block ", currentBlock.Number())
		//	return nil
		//}

		currBlock, err := s.bscClient.BlockByNumber(ctx, nil)

		if err != nil {
			return errors.Wrap(err, "failed to fetch current block")
		}

		if s.lastBlockNumber > currBlock.NumberU64() {
			s.log.Infof("Current block was already processed", currBlock.NumberU64())
			return nil
		}

		s.log.Infof("Fetching events for address", common.HexToAddress(s.contractAddress.Address))

		logs, err := s.bscClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(s.lastBlockNumber)),
			ToBlock:   currBlock.Number(),
			Addresses: []common.Address{
				common.HexToAddress(s.contractAddress.Address),
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to set filter query for filter logs")
		}

		s.lastBlockNumber = currBlock.NumberU64() + 1

		err = s.process(ctx, logs)
		if err != nil {
			return errors.Wrap(err, "failed to process events")
		}

		return nil
	}, 5*time.Second, 30*time.Second, time.Minute)
}

func (s *Service) process(ctx context.Context, logs []types.Log) error {
	s.log.Infof("Running process(), logs amount", len(logs))
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

	s.log.Infof("Started processing some received events")

	if event.Topics == nil {
		return nil
	}

	s.log.Infof("Events topic", event.Topics[0])

	if event.Topics[0].String() != s.eventsConfig.MissionCreatedHash &&
		event.Topics[0].String() != s.eventsConfig.MissionJoinedHash &&
		event.Topics[0].String() != s.eventsConfig.RewardWithdrawnHash {
		return nil
	}

	switch event.Topics[0].String() {
	case s.eventsConfig.MissionCreatedHash:
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

		s.log.Info(missionDB)
		//TODO: mission check

		missionCheck, err := s.missionQ.FilterById(missionDB.MissionId).Get()

		if err != nil {
			return errors.Wrap(err, "failed to get from db, mission")
		}

		if missionCheck == nil{
			_, err = s.missionQ.Insert(missionDB)

			if err != nil {
				return errors.Wrap(err, "failed to insert to db, mission")
			}
		}else{
			_, err = s.missionQ.Update(missionDB)

			if err != nil {
				return errors.Wrap(err, "failed to update db, mission")
			}
		}

		s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission created event")
		return nil
	case s.eventsConfig.MissionJoinedHash:
		s.log.Info("Started parsing MissionJoined event")
		missionJoined, err := spaceshipToken.ParseMissionJoined(event)
		if err != nil {
			return errors.Wrap(err, "failed parse mission joined")
		}

		investINFO, err := spaceshipToken.MissionToUsersInvest(&bind.CallOpts{}, missionJoined.MissionId, missionJoined.Player)

		if err != nil {
			return errors.Wrap(err, "failed get info about mission joined")
		}

		s.log.Infof("Invest info: ", investINFO)

		//updating mission table
		missionDB, err := s.missionQ.FilterById(0).Get()
		//missionDB, err := s.missionQ.FilterById(missionJoined.MissionId.Uint64()).Get()

		s.log.Infof("",missionDB)

		if err != nil {
			return errors.Wrap(err, "failed fetch info from db, table mission, mission")
		}

		missionFromContract, err := spaceshipToken.Missions(&bind.CallOpts{}, missionJoined.MissionId)

		if err != nil {
			return errors.Wrap(err, "failed fetch info from blockchain, mission")
		}

		if missionDB != nil {
			missionDB.MissionPower = missionFromContract.MissionPower.Int64()
			missionDB.TotalShips += investINFO.Ships.Int64()

			_, err = s.missionQ.Update(*missionDB)

			s.log.Info("Updated mission")

			if err != nil {
				return errors.Wrap(err, "failed update db, mission")
			}
		} else {
			*missionDB = data.Mission{
				MissionId:     missionJoined.MissionId.Uint64(),
				Description:   missionFromContract.Description,
				Name:          missionFromContract.Name,
				BoardingTime:  int64(missionFromContract.BoardingTime),
				LaunchTime:    int64(missionFromContract.LaunchTime),
				EndTime:       int64(missionFromContract.Duration + missionFromContract.LaunchTime),
				Duration:      int64(missionFromContract.Duration),
				MissionType:   int64(missionFromContract.MissionType),
				Reward:        missionFromContract.Reward.Int64(),
				SpaceshipCost: missionFromContract.SpaceshipCost.Int64(),
				MissionPower:  missionFromContract.MissionPower.Int64(),
				TotalShips:    investINFO.Ships.Int64(),
				NftContract:   missionFromContract.NftInfo.ContractAddress.Bytes(),
				NftTokenURI:   missionFromContract.NftInfo.TokenURI,
			}

			_, err = s.missionQ.Insert(*missionDB)

			if err != nil {
				return errors.Wrap(err, "failed insert to db, mission")
			}
		}

		//update table explorer
		explorer, err := s.explorerQ.FilterByAddress(missionJoined.Player.String()).Get()

		if err != nil {
			return errors.Wrap(err, "failed fetch info from db, table explorer, explorer")
		}

		explorerDB := data.Explorer{
			ExplorerAddress: missionJoined.Player.String(),
			TotalStakeTLM:   missionFromContract.SpaceshipCost.Int64() * investINFO.Ships.Int64(),
			TotalStakeBNB:   investINFO.BNBAmount.Int64(),
		}

		if explorer != nil {
			explorerDB.TotalStakeTLM = explorer.TotalStakeTLM + missionFromContract.SpaceshipCost.Int64() * investINFO.Ships.Int64()
			explorerDB.TotalStakeBNB = explorer.TotalStakeBNB + investINFO.BNBAmount.Int64()

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

		//updating explorer-mission table
		explorerMissionFromDB, err := s.explorer_missionQ.FilterByMission(missionJoined.MissionId.Int64()).FilterByExplorer(int64(explorerDB.ExplorerId + 1)).Get()

		//s.log.Info(s.explorer_missionQ.FilterByMission(missionJoined.MissionId.Int64()).Select())


		s.log.Info(explorerMissionFromDB)
		//s.log.Info(s.explorer_missionQ.FilterByExplorer(int64(explorerDB.ExplorerId + 1)).Get())

		if err != nil {
			return errors.Wrap(err, "failed update db, explorer-mission")
		}

		var explorerMissionDB data.ExplorerMission

		if explorerMissionFromDB == nil {
			explorerMissionDB = data.ExplorerMission{
				Explorer: int64(explorerDB.ExplorerId + 1),
				Mission:       missionJoined.MissionId.Int64(),
				Withdrawn:     false,
				NumberShips:   investINFO.Ships.Int64(),
				TotalStakeTLM: investINFO.Ships.Int64() * missionDB.SpaceshipCost,
				TotalStakeBNB: investINFO.BNBAmount.Int64(),
			}

			_, err = s.explorer_missionQ.Insert(explorerMissionDB)

			if err != nil {
				return errors.Wrap(err, "failed insert db, explorer-mission")
			}

		} else {
			explorerMissionFromDB.NumberShips += investINFO.Ships.Int64()
			explorerMissionFromDB.TotalStakeTLM = explorerMissionFromDB.NumberShips * missionDB.SpaceshipCost
			explorerMissionFromDB.TotalStakeBNB += investINFO.BNBAmount.Int64()

			_, err = s.explorer_missionQ.Update(*explorerMissionFromDB)

			if err != nil {
				return errors.Wrap(err, "failed update db, explorer-mission")
			}
		}

		s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission joined")
	case s.eventsConfig.RewardWithdrawnHash:
		rewardWithdrawn, err := spaceshipToken.ParseRewardWithdrawn(event)

		if err != nil {
			return errors.Wrap(err, "failed parse reward withdrawn")
		}

		mission, err := s.missionQ.FilterById(rewardWithdrawn.MissionId.Uint64()).Get()

		if err != nil {
			return errors.Wrap(err, "failed fetch info from db, table mission, mission")
		}

		var missionDB data.Mission

		missionToUserInvest, err := spaceshipToken.MissionToUsersInvest(&bind.CallOpts{}, rewardWithdrawn.MissionId,rewardWithdrawn.Player)

		if err != nil{
			return errors.Wrap(err, "failed fetch info from blockchain, user invest info")
		}

		if mission == nil{
			mission, err := spaceshipToken.Missions(&bind.CallOpts{}, rewardWithdrawn.MissionId)

			if err != nil {
				return errors.Wrap(err, "failed fetch info from blockchain, mission, explorer")
			}

			missionDB = data.Mission{
				MissionId:     rewardWithdrawn.MissionId.Uint64(),
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
				TotalShips:    missionToUserInvest.Ships.Int64(),
				NftContract:   mission.NftInfo.ContractAddress.Bytes(),
				NftTokenURI:   mission.NftInfo.TokenURI,
			}

			_, err = s.missionQ.Insert(missionDB)

			if err != nil {
				return errors.Wrap(err, "failed to insert to db, table mission")
			}
		}else{
			missionDB = *mission
		}

		explorer, err := s.explorerQ.FilterByAddress(rewardWithdrawn.Player.String()).Get()

		if err != nil {
			return errors.Wrap(err, "failed fetch info from db, table explorer, explorer")
		}

		var explorerDB data.Explorer

		if explorer == nil{
			explorerDB = data.Explorer{
				ExplorerAddress: rewardWithdrawn.Player.String(),
				TotalStakeTLM: missionDB.SpaceshipCost * missionToUserInvest.Ships.Int64(),
				TotalStakeBNB: missionToUserInvest.BNBAmount.Int64(),
			}

			_, err = s.explorerQ.Insert(explorerDB)

			if err != nil {
				return errors.Wrap(err, "failed to insert to db, table explorer")
			}
		}else{
			explorerDB = *explorer
		}


		explorerMission, err := s.explorer_missionQ.FilterByMission(rewardWithdrawn.MissionId.Int64()).FilterByExplorer(int64(explorerDB.ExplorerId)).Get()

		if err != nil {
			return errors.Wrap(err, "failed get info about mission-explorer")
		}

		var explorerMissionDB data.ExplorerMission

		if explorerMission == nil{
			explorerMissionDB = data.ExplorerMission{
				Explorer: int64(explorerDB.ExplorerId + 1),
				Mission: int64(missionDB.MissionId),
				Withdrawn: true,
				NumberShips: missionToUserInvest.Ships.Int64(),
				TotalStakeTLM: missionDB.SpaceshipCost * missionToUserInvest.Ships.Int64(),
				TotalStakeBNB: missionToUserInvest.BNBAmount.Int64(),
			}

			_, err = s.explorer_missionQ.Insert(explorerMissionDB)

			if err != nil {
				return errors.Wrap(err, "failed to insert to db, explorer-mission")
			}
		}else{
			explorerMissionDB = *explorerMission
			explorerMissionDB.Withdrawn = true

			_, err = s.explorer_missionQ.Update(explorerMissionDB)

			if err != nil {
				return errors.Wrap(err, "failed to update db, explorer-mission")
			}
		}

		s.log.WithField("mission_id", explorerMissionDB.Mission).Info("Success get reward withdrawn")
	default:
		return nil
	}

	return nil
}
