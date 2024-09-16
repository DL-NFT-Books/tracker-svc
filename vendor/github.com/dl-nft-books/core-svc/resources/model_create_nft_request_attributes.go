/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateNftRequestAttributes struct {
	// Id of book to get
	BookId int64 `json:"book_id"`
	// Id of network chain
	ChainId int64 `json:"chain_id"`
	// Id of the request in the marketplace contract
	MarketplaceRequestId int64 `json:"marketplace_request_id"`
	// Address of a nft collection which contains the nft to exchange
	NftAddress string `json:"nft_address"`
	// Id of NFT to exchange
	NftId int64 `json:"nft_id"`
	// Address of a user who sent the request
	Requester string `json:"requester"`
}
