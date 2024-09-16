package models

import (
	"github.com/dl-nft-books/core-svc/internal/service/api/requests"
	"github.com/dl-nft-books/core-svc/resources"
)

type (
	ListNftRequestRequest  requests.ListNftRequestsRequest
	ListNftRequestResponse resources.NftRequestListResponse
	NftRequestResponse     resources.NftRequestResponse
)
