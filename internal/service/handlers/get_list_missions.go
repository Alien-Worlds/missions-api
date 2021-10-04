package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alien-Worlds/missions-api/internal/data"
	"github.com/Alien-Worlds/missions-api/internal/service/helpers"
	"github.com/Alien-Worlds/missions-api/internal/service/requests"
	"github.com/Alien-Worlds/missions-api/resources"
	"github.com/ethereum/go-ethereum/common"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetListMissions(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetMissionListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	missionQ := helpers.Mission(r)
	missions, err := missionQ.Select(req.OffsetPageParams)

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get missions")
		ape.Render(w, problems.InternalError())
		return
	}


	result := resources.MissionListResponse{
		Data: newMissionsList(missions),
	}

	ape.Render(w, result)
}

func newMissionsList(missions []data.Mission) []resources.Mission {
	result := make([]resources.Mission, len(missions))
	for i, mission := range missions {
		result[i] = newMissionModel(mission)
	}
	return result
}

func newMissionModel(mission data.Mission) resources.Mission {
	return resources.Mission{
		Key: resources.Key{
			ID: strconv.FormatUint(mission.MissionId, 10),
			Type: resources.MISSION,
		},
		Attributes: resources.MissionAttributes{
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
		},
	}
}