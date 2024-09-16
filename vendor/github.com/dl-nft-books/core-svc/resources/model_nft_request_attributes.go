/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type NftRequestAttributes struct {
	// Id of network chain
	ChainId int64 `json:"chain_id"`
	// Timestamp when the user has created the request a book
	CreatedAt time.Time `json:"created_at"`
	// Last updated status time
	LastUpdatedAt time.Time `json:"last_updated_at"`
	// Id of the request in the marketplace contract
	MarketplaceRequestId int64 `json:"marketplace_request_id"`
	// Address of a nft collection which contains the nft to exchange
	NftAddress string `json:"nft_address"`
	// Id of NFT to exchange
	NftId int64 `json:"nft_id"`
	// Address of a user who sent the request
	Requester string `json:"requester"`
	// Nft request status
	Status NftRequestStatus `json:"status"`
}
