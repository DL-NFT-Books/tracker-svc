package ethreader

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/nft-books/contract-tracker/internal/data/ethereum"
	"gitlab.com/tokend/nft-books/contract-tracker/internal/reader"
	"gitlab.com/tokend/nft-books/contract-tracker/solidity/generated/tokenfactory"
)

type FactoryContractReader struct {
	rpc *ethclient.Client

	from    *uint64
	to      *uint64
	address *common.Address
	ctx     context.Context

	// contractInstancesCache is a map storing already initialized instances of contracts
	contractInstancesCache map[common.Address]*tokenfactory.Tokenfactory

	// rpcInstancesCache is a map storing already initialized instances of RPC connections
	rpcInstancesCache map[string]*ethclient.Client
}

func NewFactoryContractReader() reader.FactoryReader {
	return &FactoryContractReader{
		ctx:                    context.Background(),
		rpcInstancesCache:      map[string]*ethclient.Client{},
		contractInstancesCache: make(map[common.Address]*tokenfactory.Tokenfactory),
	}
}

func (r *FactoryContractReader) GetRPCInstance(rawURL string) (*ethclient.Client, error) {
	// Searching RPC instance in cache, if not found -- create new and store
	cacheInstance, ok := r.rpcInstancesCache[rawURL]
	if ok {
		return cacheInstance, nil
	}

	newInstance, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert value into eth client", logan.F{
			"raw_url": rawURL,
		})
	}

	r.rpcInstancesCache[rawURL] = newInstance
	return newInstance, nil

}

func (r *FactoryContractReader) From(from uint64) reader.FactoryReader {
	r.from = &from
	return r
}

func (r *FactoryContractReader) To(to uint64) reader.FactoryReader {
	r.to = &to
	return r
}

func (r *FactoryContractReader) WithAddress(address common.Address) reader.FactoryReader {
	r.address = &address
	return r
}

func (r *FactoryContractReader) WithCtx(ctx context.Context) reader.FactoryReader {
	r.ctx = ctx
	return r
}

func (r *FactoryContractReader) WithRPC(rpc *ethclient.Client) reader.FactoryReader {
	r.rpc = rpc
	return r
}

func (r *FactoryContractReader) validateParameters() error {
	//TODO: SHOULD WE VALIDATE `TO` PARAM?

	if r.from == nil {
		return reader.FromNotSpecifiedErr
	}
	if r.address == nil {
		return reader.AddressNotSpecifiedErr
	}
	if r.rpc == nil {
		return reader.RPCNotSpecifiedErr
	}

	return nil
}

// GetContractCreatedEvents returns the deploy contract events
// in the form of ethereum.ContractCreatedEvent array
// based on contract, start and end blocks to search through
func (r *FactoryContractReader) GetContractCreatedEvents() (events []ethereum.ContractCreatedEvent, err error) {
	if err = r.validateParameters(); err != nil {
		return nil, err
	}

	instance, err := r.getInstance(*r.address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get instance")
	}

	iterator, err := instance.FilterTokenContractDeployed(
		&bind.FilterOpts{
			Start: *r.from,
			End:   r.to,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize an iterator")
	}
	if iterator == nil {
		return nil, errors.From(NullIteratorErr, logan.F{
			"contract": r.address.String(),
		})
	}

	defer func(iterator *tokenfactory.TokenfactoryTokenContractDeployedIterator) {
		if tempErr := iterator.Close(); tempErr != nil {
			err = tempErr
		}
	}(iterator)

	for iterator.Next() {
		event := iterator.Event
		if event != nil {
			receipt, err := r.rpc.TransactionReceipt(r.ctx, event.Raw.TxHash)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
					"tx_hash": event.Raw.TxHash.String(),
				})
			}

			events = append(events, ethereum.ContractCreatedEvent{
				Address:     event.NewTokenContractAddr,
				BlockNumber: event.Raw.BlockNumber,
				Name:        event.TokenName,
				Symbol:      event.TokenSymbol,
				Status:      receipt.Status,
			})
		}
	}

	return events, nil
}

func (r *FactoryContractReader) getInstance(address common.Address) (*tokenfactory.Tokenfactory, error) {
	cacheInstance, ok := r.contractInstancesCache[address]
	if ok {
		return cacheInstance, nil
	}

	newInstance, err := tokenfactory.NewTokenfactory(address, r.rpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize token factory instance for given address", logan.F{
			"address": address,
		})
	}

	r.contractInstancesCache[address] = newInstance
	return newInstance, nil
}
