package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type SignMintByNftRequest struct {
	TaskID     int64  `url:"task_id"`
	NftAddress string `url:"token_address"`
}

func NewSignMintByNftRequest(r *http.Request) (*SignMintByNftRequest, error) {
	var result SignMintByNftRequest
	var err error

	if err = urlval.Decode(r.URL.Query(), &result); err != nil {
		return &result, errors.Wrap(err, "failed to unmarshal sign mint request")
	}

	return &result, result.validate()
}

func (r SignMintByNftRequest) validate() error {
	return validation.Errors{
		"task_id=": validation.Validate(
			r.TaskID,
			validation.Required,
			validation.Min(1)),
	}.Filter()

}
