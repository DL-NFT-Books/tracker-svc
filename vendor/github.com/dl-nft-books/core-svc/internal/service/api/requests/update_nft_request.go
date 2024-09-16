package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"

	"github.com/dl-nft-books/core-svc/resources"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type UpdateNftRequestRequest struct {
	ID int64
	resources.UpdateNftRequestRequest
}

func NewUpdateNftRequestRequest(r *http.Request) (*UpdateNftRequestRequest, error) {
	var request UpdateNftRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal update task request")
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id from the url path")
	}
	request.ID = int64(id)

	return &request, request.validate()
}

func (r UpdateNftRequestRequest) validate() error {
	return validation.Errors{
		"data/attributes/status": validation.Validate(
			r.Data.Attributes.Status,
			validation.Min(1),
			validation.Max(5),
			validation.Required,
		),
	}.Filter()
}
