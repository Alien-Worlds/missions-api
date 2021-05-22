package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type MissionListRequest struct {
	pgdb.OffsetPageParams
}

func NewGetMissionListRequest(r *http.Request) (MissionListRequest, error) {
	request := MissionListRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}