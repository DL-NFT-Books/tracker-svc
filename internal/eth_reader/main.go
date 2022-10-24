package eth_reader

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/nft-books/contract-tracker/solidity/generated/itokencontract"
)

var NullIteratorErr = errors.New("iterator has a nil value")

type TokenContractReader struct {
	rpc *ethclient.Client
}

func NewTokenContractReader(rpc *ethclient.Client) TokenContractReader {
	return TokenContractReader{
		rpc: rpc,
	}
}

type TokenMintedEvent struct {
	Recipient common.Address
	TokenId   int64
	Uri       string
}

func (r *TokenContractReader) GetEvents(
	contract common.Address,
	startBlock,
	endBlock uint64,
) (
	events []TokenMintedEvent,
	lastBlock uint64,
	err error,
) {
	instance, err := itokencontract.NewItokencontract(contract, r.rpc)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to initialize a contract instance")
	}

	iterator, err := instance.FilterTokenMinted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   &endBlock,
		}, nil,
	)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to initialize an iterator")
	}
	if iterator == nil {
		return nil, 0, NullIteratorErr
	}

	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event

		if event != nil {
			events = append(events, TokenMintedEvent{
				Recipient: event.Recipient,
				TokenId:   event.TokenId.Int64(),
				Uri:       event.TokenURI,
			})

			lastBlock = event.Raw.BlockNumber
		}
	}

	return
}
