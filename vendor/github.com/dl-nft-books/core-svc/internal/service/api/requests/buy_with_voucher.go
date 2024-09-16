package requests

import (
	"encoding/json"
	"github.com/dl-nft-books/core-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"net/http"
)

type BuyWithVoucherRequest struct {
	resources.BuyWithVoucherRequest
}

func NewBuyWithVoucherRequest(r *http.Request) (*BuyWithVoucherRequest, error) {
	var request BuyWithVoucherRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal request")
	}

	return &request, request.validate()
}

func (r BuyWithVoucherRequest) validate() error {
	return validation.Errors{
		"task_id": validation.Validate(
			r.Data.Attributes.TaskId,
			validation.Required,
			validation.Min(1)),
		"sig":             validation.Validate(r.Data.Attributes.PermitSig, validation.Required),
		"voucher_address": validation.Validate(r.Data.Attributes.VoucherAddress, validation.Required),
	}.Filter()
}
