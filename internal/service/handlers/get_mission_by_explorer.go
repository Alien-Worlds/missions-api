package handlers

import (
	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/Alien-Worlds/missions-api/internal/service/helpers"
	"github.com/Alien-Worlds/missions-api/resources"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetMissionsByExplorerAddress(w http.ResponseWriter, r *http.Request) {
	explorerAddressString := chi.URLParam(r, "explorer-address")

	explorerMissionQ := helpers.ExplorerMission(r)

	explorerQ := helpers.Explorer(r)
	explorer, err := explorerQ.FilterByAddress(explorerAddressString).Get()

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get explorer from db")
		ape.Render(w, problems.InternalError())
		return
	}

	if explorer == nil{
		helpers.Log(r).WithError(err).Error("not found explorer from db")
		ape.Render(w, problems.NotFound())
		return
	}

	explorerMissions, err := explorerMissionQ.FilterByExplorer(int64(explorer.ExplorerId)).Select()

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get explorer-missions")
		ape.Render(w, problems.InternalError())
		return
	}

	missionQ := helpers.Mission(r)
	missionByExplorerMissionsList := make([]resources.MissionByExplorerMissions, len(explorerMissions))

	for i, explorerMission := range explorerMissions {
		mission, err := missionQ.FilterById(explorerMission.Mission).Get()

		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get mission")
			ape.Render(w, problems.InternalError())
			return
		}

		missionByExplorerMissionsList[i] = newMissionByExplorerMissionsModel(*mission, explorerMission)
	}

	result := resources.MissionByExplorerResponse{
		Data: resources.MissionByExplorer{
			Key: resources.Key{
				ID:   strconv.FormatUint(explorer.ExplorerId, 10),
				Type: resources.EXPLORER,
			},
			Attributes: resources.MissionByExplorerAttributes{
				Address:  explorer.ExplorerAddress,
				Missions: missionByExplorerMissionsList,
				TotalInvestInfo: resources.MissionByExplorerInvestInfo{
					TotalStakeBNB: explorer.TotalStakeBNB,
					TotalStakeTLM: explorer.TotalStakeTLM,
				},
			},
		},
	}

	ape.Render(w, result)
}

func newMissionByExplorerMissionsModel(mission data.Mission, explorerMission data.ExplorerMission) resources.MissionByExplorerMissions {
	return resources.MissionByExplorerMissions{
		Key: resources.Key{
			ID:   strconv.FormatUint(mission.MissionId, 10),
			Type: resources.MISSION,
		},
		Attributes: resources.MissionByExplorerMissionsAttributes{
			BoardingTime:  mission.BoardingTime,
			Description:   mission.Description,
			Duration:      mission.Duration,
			EndTime:       mission.EndTime,
			LaunchTime:    mission.LaunchTime,
			MissionPower:  mission.MissionPower,
			MissionType:   mission.MissionType,
			Name:          mission.Name,
			NftContract:   common.BytesToAddress(mission.NftContract).String(),
			NftTokenURI:   mission.NftTokenURI,
			Reward:        mission.Reward,
			SpaceshipCost: mission.SpaceshipCost,
			TotalShips:    mission.TotalShips,
			InvestInfo: resources.InvestInfo{
				NumberOfShips: explorerMission.NumberShips,
				TotalStakeBNB: explorerMission.TotalStakeBNB,
				TotalStakeTLM: explorerMission.TotalStakeTLM,
				Withdrawn:     explorerMission.Withdrawn,
			},
		},
	}
}
