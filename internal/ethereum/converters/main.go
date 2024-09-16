package converters

import (
	"context"
	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"github.com/dl-nft-books/tracker-svc/resources"
	"github.com/dl-nft-books/tracker-svc/solidity/generated/erc20"
	"github.com/dl-nft-books/tracker-svc/solidity/generated/marketplace"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
	"time"
)

type EventConverter struct {
	client      *ethclient.Client
	ctx         context.Context
	nativeToken etherdata.Erc20Info
}

func NewEventConverter(client *ethclient.Client, ctx context.Context, nativeToken etherdata.Erc20Info) EventConverter {
	return EventConverter{
		client:      client,
		ctx:         ctx,
		nativeToken: nativeToken,
	}
}

func (c *EventConverter) TokenSuccessfullyPurchased(raw marketplace.MarketplaceTokenSuccessfullyPurchased) (*etherdata.TokenSuccessfullyPurchasedEvent, error) {
	receipt, err := c.client.TransactionReceipt(c.ctx, raw.Raw.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
			"tx_hash": raw.Raw.TxHash.String(),
		})
	}

	purchaseTimestamp, err := c.getBlockTimestamp(raw.Raw.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block timestamp")
	}

	erc20Data, err := c.GetErc20Data(raw.BuyParams.PaymentDetails.PaymentTokenAddress, raw.PaymentType)
	event := &etherdata.TokenSuccessfullyPurchasedEvent{
		ContractAddress:   raw.BuyParams.TokenContract,
		Recipient:         raw.BuyParams.Recipient,
		TokenId:           raw.BuyParams.TokenData.TokenId.Int64(),
		NftId:             raw.BuyParams.PaymentDetails.NftTokenId.Int64(),
		Uri:               raw.BuyParams.TokenData.TokenURI,
		MintedTokenPrice:  raw.MintedTokenPrice,
		Amount:            raw.PaidTokensAmount,
		Erc20Info:         *erc20Data,
		PaymentTokenPrice: raw.BuyParams.PaymentDetails.PaymentTokenPrice,
		Status:            receipt.Status,
		BlockNumber:       raw.Raw.BlockNumber,
		Timestamp:         *purchaseTimestamp,
		Type:              raw.PaymentType,
	}

	return event, nil
}

func (c *EventConverter) NFTRequestCreated(raw marketplace.MarketplaceNFTRequestCreated) (*etherdata.NFTRequestCreatedEvent, error) {
	receipt, err := c.client.TransactionReceipt(c.ctx, raw.Raw.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
			"tx_hash": raw.Raw.TxHash.String(),
		})
	}

	requestTimestamp, err := c.getBlockTimestamp(raw.Raw.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block timestamp")
	}

	event := &etherdata.NFTRequestCreatedEvent{
		RequestId:     raw.RequestId,
		Requester:     raw.Requester,
		TokenContract: raw.TokenContract,
		NftContract:   raw.NftContract,
		NftId:         raw.NftId,
		Status:        receipt.Status,
		BlockNumber:   raw.Raw.BlockNumber,
		Timestamp:     *requestTimestamp,
	}

	return event, nil
}

func (c *EventConverter) NFTRequestCanceled(raw marketplace.MarketplaceNFTRequestCanceled) (*etherdata.NFTRequestCanceledEvent, error) {
	receipt, err := c.client.TransactionReceipt(c.ctx, raw.Raw.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
			"tx_hash": raw.Raw.TxHash.String(),
		})
	}

	purchaseTimestamp, err := c.getBlockTimestamp(raw.Raw.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block timestamp")
	}

	event := &etherdata.NFTRequestCanceledEvent{
		RequestId:   raw.RequestId,
		Status:      receipt.Status,
		BlockNumber: raw.Raw.BlockNumber,
		Timestamp:   *purchaseTimestamp,
	}

	return event, nil
}

func (c *EventConverter) TokenSuccessfullyExchanged(raw marketplace.MarketplaceTokenSuccessfullyExchanged) (*etherdata.TokenSuccessfullyExchangedEvent, error) {
	receipt, err := c.client.TransactionReceipt(c.ctx, raw.Raw.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt", logan.F{
			"tx_hash": raw.Raw.TxHash.String(),
		})
	}

	purchaseTimestamp, err := c.getBlockTimestamp(raw.Raw.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block timestamp")
	}

	event := &etherdata.TokenSuccessfullyExchangedEvent{
		TokenContract: raw.NftRequestInfo.TokenContract,
		Recipient:     raw.AcceptRequestParams.Recipient,
		RequestId:     raw.AcceptRequestParams.RequestId,
		TokenId:       raw.AcceptRequestParams.TokenData.TokenId,
		NftAddress:    raw.NftRequestInfo.NftContract,
		NftId:         raw.NftRequestInfo.NftId,
		Uri:           raw.AcceptRequestParams.TokenData.TokenURI,
		Status:        receipt.Status,
		BlockNumber:   raw.Raw.BlockNumber,
		Timestamp:     *purchaseTimestamp,
	}

	return event, nil
}

// getBlockTimestamp is a function that returns a timestamp
// when a block with specified blockNumber was initialized
func (c *EventConverter) getBlockTimestamp(blockNumber uint64) (*time.Time, error) {
	// Get header by a block number
	header, err := c.client.BlockByNumber(c.ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block header by its number")
	}

	// Get timestamp from a block as an uint64 and convert it to time.Time
	blockTime := time.Unix(int64(header.Time()), 0)
	return &blockTime, nil
}

// GetErc20Data is a function that retrieves the etherdata.Erc20Info from the
// address specified in parameters. If address is etherdata.NullAddress, function
// returns the native erc20 data. Otherwise, contract on the specified
// address must implement IERC20 interface.
func (c *EventConverter) GetErc20Data(address common.Address, eventType uint8) (*etherdata.Erc20Info, error) {
	if resources.TokenPurchasedEventType(eventType) == resources.NFT {
		return &etherdata.Erc20Info{
			TokenAddress: address,
			Name:         "",
			Symbol:       "NFT",
			Decimals:     0,
		}, nil
	}
	if address == etherdata.NullAddress {
		return &c.nativeToken, nil
	}

	erc20Instance, err := erc20.NewErc20(address, c.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize erc20 instance", logan.F{
			"erc20_address": address,
		})
	}

	tokenName, err := erc20Instance.Name(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read marketplace's name from the contract instance")
	}

	tokenSymbol, err := erc20Instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read marketplace's symbol from the contract instance")
	}

	tokenDecimals, err := erc20Instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get marketplace's decimals")
	}

	return &etherdata.Erc20Info{
		TokenAddress: address,
		Name:         tokenName,
		Symbol:       tokenSymbol,
		Decimals:     tokenDecimals,
	}, nil
}
