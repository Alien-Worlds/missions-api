package handlers

import (
	"github.com/Alien-Worlds/missions-api/internal/service/helpers"
	"github.com/Alien-Worlds/missions-api/resources"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetMissionById(w http.ResponseWriter, r *http.Request) {
	missionIdString := chi.URLParam(r, "mission-id")

	missionQ := helpers.Mission(r)

	missionId, err := strconv.ParseInt(missionIdString, 10, 64)

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse mission id")
		ape.Render(w, problems.InternalError())
		return
	}

	mission, err := missionQ.FilterById(missionId).Get()

	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get mission")
		ape.Render(w, problems.InternalError())
		return
	}

	if mission == nil{
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.MissionResponse{
		Data: newMissionModel(*mission),
	}

	ape.Render(w, result)
}

