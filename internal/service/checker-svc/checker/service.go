package checker

import (
	"context"
	eth "github.com/ethereum/go-ethereum"
	bind2 "github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"

	//_ "github.com/binance-chain/bsc-static/bsc"
	"github.com/binance-chain/bsc-static/bsc/accounts/abi/bind"
	//_ "github.com/binance-chain/bsc-static/bsc/common"
	"github.com/binance-chain/bsc-static/bsc/core/types"
	//_ "github.com/lib/pq"
	"github.com/redcuckoo/bsc-checker-events/internal/contracts"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Running checker service")

	running.WithBackOff(ctx, s.log, "new-checker-service", func(ctx context.Context) error {
		var ch chan<- *ethTypes.Header

		if s.lastBlockNumber == 0 {
			currentBlock, err := s.bscClient.BlockByNumber(ctx, nil)

			if err != nil {
				return errors.Wrap(err, "startup error, failed to get last block")
			}

			s.lastBlockNumber = currentBlock.Header().Number.Uint64()
		}

		sub, err := s.bscClient.SubscribeNewHead(ctx, ch)

		if err != nil {
			return errors.Wrap(err, "subscription on blockchain head failed")
		}

		for c := range ch {
			logs, err := s.bscClient.FilterLogs(ctx, eth.FilterQuery{
				FromBlock: big.NewInt(int64(c.Number.Uint64() + 1)),
				Addresses: []ethCommon.Address{
					ethCommon.Address(s.contractAddress),
				},
			})

			if err != nil {
				return errors.Wrap(err, "failed to set filter query for filter logs")
			}
			err = s.process(ctx, logs)
			if err != nil {
				return errors.Wrap(err, "failed to process events")
			}

		}

		sub.Unsubscribe()
		return nil
	}, 5*time.Second, 30*time.Second, time.Minute)
}

func (s *Service) process(ctx context.Context, logs []ethTypes.Log) error {
	for _, event := range logs {
		err := s.processTransfer(ctx, types.Log(event))
		if err != nil {
			s.log.WithError(err).Error("failed to process transfer")
			continue
		}
	}

	return nil
}

func (s *Service) processTransfer(ctx context.Context, event types.Log) error {
	spaceshipToken, err := SpaceshipStaking.NewSpaceshipStaking(ethCommon.Address(s.contractAddress), &s.bscClient)
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
		missionCreated, err := spaceshipToken.ParseMissionCreated(ethTypes.Log(event))
		if err != nil {
			return errors.Wrap(err, "failed parse mission created")
		}
		mission, err := spaceshipToken.Missions((*bind2.CallOpts)(&bind.CallOpts{}), missionCreated.Id)
		if err != nil {
			return errors.Wrap(err, "failed get info about mission created")
		}
		missionDB := data.Mission{
			MissionId:     missionCreated.Id.Uint64(),
			Description:   mission.Description,
			Name:          mission.Name,
			BoardingTime:  mission.BoardingTime,
			LaunchTime:    mission.LaunchTime,
			EndTime:       mission.Duration + mission.LaunchTime,
			Duration:      mission.Duration,
			MissionType:   mission.MissionType,
			Reward:        mission.Reward.Uint64(),
			SpaceshipCost: mission.SpaceshipCost.Uint64(),
			MissionPower:  mission.MissionPower.Uint64(),
			TotalShips:    0,
			NftContract:   mission.NftInfo.ContractAddress.Bytes(),
			NftTokenURI:   mission.NftInfo.TokenURI,
		}

		//TODO: mission check
		_, err = s.missionQ.Insert(missionDB)
		if err != nil {
			return errors.Wrap(err, "failed to insert to db")
		}

		s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission")
	case s.eventsConfig.MissionJoinedHash:
		missionJoined, err := spaceshipToken.ParseMissionJoined(ethTypes.Log(event))
		if err != nil {
			return errors.Wrap(err, "failed parse mission joined")
		}
		mission, err := spaceshipToken.Missions((*bind2.CallOpts)(&bind.CallOpts{}), missionJoined.MissionId)

		if err != nil {
			return errors.Wrap(err, "failed get info about mission joined")
		}

		investINFO, err := spaceshipToken.MissionToUsersInvest((*bind2.CallOpts)(&bind.CallOpts{}),missionJoined.MissionId, missionJoined.Player)

		explorerMissionDB := data.ExplorerMission{
			Explorer:      missionJoined.Player.String(),
			Mission:       missionJoined.MissionId.Uint64(),
			Withdrawn:     false,
			NumberShips:   investINFO.Ships.Uint64(),
			TotalStakeTLM: investINFO.Ships.Uint64() * mission.SpaceshipCost.Uint64(),
			TotalStakeBNB: investINFO.BNBAmount.Uint64(),
		}

		totalStakedTLM, err := s.explorerQ.SelectTotalTLM(missionJoined.Player.String())
		totalStakedBNB, err := s.explorerQ.SelectTotalBNB(missionJoined.Player.String())

		explorerDB := data.Explorer{
			ExplorerId:    missionJoined.Player.String(),
			TotalStakeTLM: totalStakedTLM.TotalStakeTLM + explorerMissionDB.TotalStakeTLM,
			TotalStakeBNB: totalStakedBNB.TotalStakeBNB + investINFO.BNBAmount.Uint64(),
		}

		missionDB := data.Mission{
			MissionId:     missionJoined.MissionId.Uint64(),
			Description:   mission.Description,
			Name:          mission.Name,
			BoardingTime:  mission.BoardingTime,
			LaunchTime:    mission.LaunchTime,
			EndTime:       mission.Duration + mission.LaunchTime,
			Duration:      mission.Duration,
			MissionType:   mission.MissionType,
			Reward:        mission.Reward.Uint64(),
			SpaceshipCost: mission.SpaceshipCost.Uint64(),
			MissionPower:  mission.MissionPower.Uint64(),
			TotalShips:    investINFO.Ships.Uint64(),
			NftContract:   mission.NftInfo.ContractAddress.Bytes(),
			NftTokenURI:   mission.NftInfo.TokenURI,
		}

		//TODO: mission check
		_, err = s.missionQ.Update(missionDB)
		_, err = s.explorerQ.Update(explorerDB)
		_, err = s.explorer_missionQ.Insert(explorerMissionDB)
		if err != nil {
			return errors.Wrap(err, "failed to insert/update to db")
		}

		s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission joined")
	case s.eventsConfig.RewardWithdrawnHash:
		rewardWithdrawn, err := spaceshipToken.ParseRewardWithdrawn(ethTypes.Log(event))
		if err != nil {
			return errors.Wrap(err, "failed parse reward withdrawn")
		}
		mission, err := spaceshipToken.Missions((*bind2.CallOpts)(&bind.CallOpts{}), rewardWithdrawn.MissionId)
		if err != nil {
			return errors.Wrap(err, "failed get info about mission ")
		}
		missionDB := data.Mission{
			MissionId:     rewardWithdrawn.MissionId.Uint64(),
			Description:   mission.Description,
			Name:          mission.Name,
			BoardingTime:  mission.BoardingTime,
			LaunchTime:    mission.LaunchTime,
			EndTime:       mission.Duration + mission.LaunchTime,
			Duration:      mission.Duration,
			MissionType:   mission.MissionType,
			Reward:        mission.Reward.Uint64(),
			SpaceshipCost: mission.SpaceshipCost.Uint64(),
			MissionPower:  mission.MissionPower.Uint64(),
			TotalShips:    0,
			NftContract:   mission.NftInfo.ContractAddress.Bytes(),
			NftTokenURI:   mission.NftInfo.TokenURI,
		}

		//TODO: mission check
		_, err = s.missionQ.Insert(missionDB)
		if err != nil {
			return errors.Wrap(err, "failed to insert to db")
		}

		s.log.WithField("mission_id", missionDB.MissionId).Info("Success get mission")
	default:
		return nil
	}

	//redeemed, err := spaceshipToken.ParseMissionCreated(ethTypes.Log(event))
	//if err != nil {
	//	return errors.Wrap(err, "failed parse redemption requested")
	//}

	//redemption, err := spaceshipToken.Missions((*bind2.CallOpts)(&bind.CallOpts{}), redeemed.Id)
	//if err != nil {
	//	return errors.Wrap(err, "failed get info about redemption requests")
	//}

	//if redemption.Duration == 0 {
	//	return nil
	//}

	//decimals, err := spaceshipToken.Decimals(&bind.CallOpts{})
	//if err != nil {
	//	return errors.Wrap(err, "failed to get decimals from contracts")
	//}
	//
	//amount := s.prepareAmount(7, int64(decimals), redemption.Amount)

	//withdraw := data.Withdraw{
	//	Amount:        amount,
	//	StellarWallet: redemption.Recipient,
	//	TxHashEth:     event.TxHash.String(),
	//	State:         data.StatesPending,
	//	AssetCode:     s.asset.AssetCode,
	//	RedemptionId:  redeemed.RedemptionRequestID.String(),
	//	NumberBlock:   int64(event.BlockNumber),
	//	Signers:       pq.StringArray{},
	//}



	//conflict, err := s.withdrawQ.FilterByTxHashEth(event.TxHash.String()).Get()
	//if err != nil {
	//	return errors.Wrap(err, "failed to get withdraw by eth hash")
	//}
	//
	//if conflict.TxHashEth == event.TxHash.String() {
	//	return nil
	//}

	//_, err = s.missionQ.Insert(mission)
	//if err != nil {
	//	return errors.Wrap(err, "failed to insert to db")
	//}
	//
	//s.log.WithField("redemption_id", mission.MissionId).Info("Success get mission")

	return nil
}
