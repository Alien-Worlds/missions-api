package handlers

import (
	"github.com/go-chi/chi"
	"github.com/redcuckoo/bsc-checker-events/internal/service/helpers"
	"github.com/redcuckoo/bsc-checker-events/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetMissionsByExplorerId(w http.ResponseWriter, r *http.Request) {
	explorerIdString := chi.URLParam(r, "explorer-id")

	explorerMissionQ := helpers.ExplorerMission(r)

	explorerId, err := strconv.ParseInt(explorerIdString, 10, 64)

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse explorer id")
		ape.Render(w, problems.InternalError())
		return
	}

	explorerMissions, err := explorerMissionQ.FilterByExplorer(explorerId).Select()

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get explorer-missions")
		ape.Render(w, problems.InternalError())
		return
	}

	missionQ := helpers.Mission(r)
	missionList := make([]resources.Mission, len(explorerMissions))

	for i, explorerMission := range explorerMissions {
		mission, err := missionQ.FilterById(explorerMission.Mission).Get()

		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get mission")
			ape.Render(w, problems.InternalError())
			return
		}

		missionList[i] = newMissionModel(*mission)
	}

	result := resources.MissionListResponse{
		Data: missionList,
	}

	ape.Render(w, result)
}
