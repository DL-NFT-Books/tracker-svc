package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"gitlab.com/tokend/nft-books/generator-svc/resources"
)

type ListTokensRequest struct {
	pgdb.OffsetPageParams
	Sorts pgdb.Sorts `url:"sort" default:"id"`

	Account      []string                `filter:"account"`
	Status       []resources.TokenStatus `filter:"status"`
	TokenId      *int64                  `filter:"token_id"`
	MetadataHash []string                `filter:"metadata_hash"`
}

func NewListTokensRequest(r *http.Request) (*ListTokensRequest, error) {
	var request ListTokensRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
