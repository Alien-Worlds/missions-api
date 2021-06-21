package handlers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/redcuckoo/bsc-checker-events/internal/data"
	"github.com/redcuckoo/bsc-checker-events/internal/service/helpers"
	"github.com/redcuckoo/bsc-checker-events/resources"
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

	if err != nil || explorer == nil {
		helpers.Log(r).WithError(err).Error("failed to get explorer from db")
		ape.Render(w, problems.InternalError())
		return
	}

	explorerMissions, err := explorerMissionQ.FilterByExplorer(int64(explorer.ExplorerId)).Select()

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get explorer-missions")
		ape.Render(w, problems.InternalError())
		return
	}

	missionQ := helpers.Mission(r)
	missionByExplorerList := make([]resources.MissionByExplorer, len(explorerMissions))

	for i, explorerMission := range explorerMissions {
		mission, err := missionQ.FilterById(explorerMission.Mission).Get()

		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get mission")
			ape.Render(w, problems.InternalError())
			return
		}

		missionByExplorerList[i] = newMissionByExplorerModel(*mission, explorerMission.Withdrawn)
	}

	result := resources.MissionByExplorerListResponse{
		Data: missionByExplorerList,
	}

	ape.Render(w, result)
}

func newMissionByExplorerModel(mission data.Mission, withdrawn bool) resources.MissionByExplorer{
	return resources.MissionByExplorer{
		Key: resources.Key{
			ID: strconv.FormatUint(mission.MissionId, 10),
			Type: resources.MISSION,
		},
		Attributes: resources.MissionByExplorerAttributes{
			BoardingTime: mission.BoardingTime,
			Description:  mission.Description,
			Duration:     mission.Duration,
			EndTime:      mission.EndTime,
			LaunchTime:   mission.LaunchTime,
			MissionPower: mission.MissionPower,
			MissionType:  mission.MissionType,
			Name:         mission.Name,
			NftContract: common.BytesToAddress(mission.NftContract).String(),
			NftTokenURI: mission.NftTokenURI,
			Reward: mission.Reward,
			SpaceshipCost: mission.SpaceshipCost,
			TotalShips:	mission.TotalShips,
			Withdrawn: withdrawn,
		},
	}
}