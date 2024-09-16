package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type ValidatePromocodeRequest struct {
	Promocode string
	Bookid    int64 `url:"book_id"`
}

func NewValidatePromocodeRequest(r *http.Request) (*ValidatePromocodeRequest, error) {
	var request ValidatePromocodeRequest
	request.Promocode = chi.URLParam(r, "promocode")

	if err := urlval.Decode(r.URL.Query(), &request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list promocodes request")
	}
	return &request, request.validate()
}

func (r *ValidatePromocodeRequest) validate() error {
	return validation.Errors{
		"book_id": validation.Validate(r.Bookid, validation.Required),
	}.Filter()
}
