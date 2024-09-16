package requests

import (
	"encoding/json"
	"net/http"

	"github.com/dl-nft-books/core-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateNftRequestRequest struct {
	resources.CreateNftRequestRequest
}

func NewCreateNftRequestRequest(r *http.Request) (*CreateNftRequestRequest, error) {
	var request CreateNftRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal create nftRequest request")
	}

	return &request, request.validate()
}

func (r CreateNftRequestRequest) validate() error {
	return validation.Errors{
		"data/attributes/": validation.Validate(
			&r.Data.Attributes,
			validation.Required,
		),
		"data/attributes/chain_id": validation.Validate(
			&r.Data.Attributes.ChainId,
			validation.Required,
		),
		"data/attributes/nft_id": validation.Validate(
			&r.Data.Attributes.NftId,
			validation.Required,
		),
		"data/attributes/nft_address": validation.Validate(
			&r.Data.Attributes.NftAddress,
			validation.Required,
		),
		"data/attributes/marketplace_request_id": validation.Validate(
			&r.Data.Attributes.MarketplaceRequestId,
			validation.Required,
		),
		"data/attributes/requester": validation.Validate(
			&r.Data.Attributes.Requester,
			validation.Required,
		),
		"data/attributes/book_id": validation.Validate(
			&r.Data.Attributes.BookId,
			validation.Required,
		),
	}.Filter()
}
