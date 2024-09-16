package etherdata

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// NFTRequestCreatedEvent is a struct containing all the data from
// the successful nft request created event
type NFTRequestCreatedEvent struct {
	RequestId     *big.Int
	Requester     common.Address
	TokenContract common.Address
	NftContract   common.Address
	NftId         *big.Int
	Status        uint64 // 0 if failed and 1 if successful
	BlockNumber   uint64
	Timestamp     time.Time
}
