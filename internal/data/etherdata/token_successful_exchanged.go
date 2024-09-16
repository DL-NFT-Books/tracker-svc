package etherdata

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// TokenSuccessfullyExchangedEvent is a struct containing all the data from
// the successful mint event which occurs after the user has
// successfully passed signature from the backend and minted a book
type TokenSuccessfullyExchangedEvent struct {
	TokenContract common.Address // address of book
	Recipient     common.Address // Who obtained a erc-721 marketplace
	TokenId       *big.Int       // Token id assigned to the erc-721 marketplace
	Uri           string         // Hash of a metadata file from the event
	Requester     common.Address // The address of the offerer
	NftAddress    common.Address // The address of the NFT contract
	NftId         *big.Int       // The ID of the NFT
	RequestId     *big.Int       // Token id assigned to the erc-721 marketplace
	Status        uint64         // 0 if failed and 2 if successful
	BlockNumber   uint64
	Timestamp     time.Time
}
