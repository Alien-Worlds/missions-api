package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/redcuckoo/bsc-checker-events/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewGetMissionRequest(r *http.Request) (resources.MissionResponse, error) {
	var request resources.MissionResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, ValidateGetMissionRequest(request)
}

func ValidateGetMissionRequest(r resources.MissionResponse) error {
	errs := validation.Errors{
		"/data/":                      validation.Validate(r.Data, validation.Required),
		"/data/data/attributes": validation.Validate(r.Data.ID),
	}

	return errs.Filter()
}