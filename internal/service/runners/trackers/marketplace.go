package trackers

import (
	"context"
	"github.com/dl-nft-books/network-svc/connector/models"
	"sync"

	"github.com/dl-nft-books/tracker-svc/internal/config"
	"github.com/dl-nft-books/tracker-svc/internal/data"
	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"github.com/dl-nft-books/tracker-svc/internal/data/postgres"
	"github.com/dl-nft-books/tracker-svc/internal/ethereum"
	"github.com/dl-nft-books/tracker-svc/internal/ethereum/listeners/marketplace"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	tokenSuccessfullyPurchasedSuffix = "-marketplace-successfully-purchased"
	tokenSuccessfullyExchangedSuffix = "-marketplace-successfully-exchanged"
	nftRequestCreatedSuffix          = "-marketplace-nft-req-created"
	nftRequestCanceledSuffix         = "-marketplace-nft-req-canceled"
)

type DeployedToken struct {
	Network models.NetworkDetailedResponse
	Address common.Address
}

type MarketPlaceTracker struct {
	log                *logan.Entry
	cfg                config.Trackers
	ctx                context.Context
	database           data.DB
	listener           ethereum.TokenListener
	marketplaceAddress common.Address
	network            models.NetworkDetailedResponse
}

func NewMarketPlaceTracker(cfg config.Config, ctx context.Context, deployedToken DeployedToken) *MarketPlaceTracker {
	mutex := new(sync.RWMutex)

	return &MarketPlaceTracker{
		log:                cfg.Log(),
		ctx:                ctx,
		cfg:                cfg.Trackers(),
		database:           postgres.NewDB(cfg.DB()),
		network:            deployedToken.Network,
		marketplaceAddress: deployedToken.Address,
		listener: token_listeners.
			NewTokenListener(ctx, mutex, deployedToken.Network).
			WithMaxDepth(cfg.Trackers().MaxDepth).
			WithDelayBetweenIntervals(cfg.Trackers().DelayBetweenIntervals),
	}
}

func (t *MarketPlaceTracker) TrackTokenSuccessfullyPurchasedEvents(ch chan<- etherdata.TokenSuccessfullyPurchasedEvent) {
	listener := t.listener.WithAddress(t.marketplaceAddress)
	running.WithBackOff(
		t.ctx,
		t.log,
		t.cfg.Prefix+tokenSuccessfullyPurchasedSuffix,
		func(ctx context.Context) error {
			startBlock := uint64(t.network.FirstBlock)
			block, err := t.database.Blocks().FilterByContractAddress(t.marketplaceAddress.String()).FilterByChainId(t.network.ChainId).Get()
			if err != nil {
				return errors.Wrap(err, "failed to get block to begin with")
			}
			if block != nil {
				startBlock = block.TokenPurchasedBlock
			}

			t.log.WithFields(logan.F{
				"marketplace_address": t.marketplaceAddress,
				"network":             t.network.Name,
			}).Info("start tracking mint events from block ", startBlock)
			return listener.
				From(startBlock).
				WithCtx(ctx).
				WatchTokenSuccessfullyPurchasedEvents(ch)
		},
		t.cfg.Backoff.NormalPeriod,
		t.cfg.Backoff.MinAbnormalPeriod,
		t.cfg.Backoff.MaxAbnormalPeriod,
	)
}

func (t *MarketPlaceTracker) TrackTokenSuccessfullyExchangedEvents(ch chan<- etherdata.TokenSuccessfullyExchangedEvent) {
	listener := t.listener.WithAddress(t.marketplaceAddress)
	running.WithBackOff(
		t.ctx,
		t.log,
		t.cfg.Prefix+tokenSuccessfullyExchangedSuffix,
		func(ctx context.Context) error {
			startBlock := uint64(t.network.FirstBlock)
			block, err := t.database.Blocks().FilterByContractAddress(t.marketplaceAddress.String()).FilterByChainId(t.network.ChainId).Get()
			if err != nil {
				return errors.Wrap(err, "failed to get block to begin with")
			}
			if block != nil {
				startBlock = block.TokenExchangedBlock
			}

			t.log.WithFields(logan.F{
				"marketplace_address": t.marketplaceAddress,
				"network":             t.network.Name,
			}).Info("start tracking mint events from block ", startBlock)
			return listener.
				From(startBlock).
				WithCtx(ctx).
				WatchTokenSuccessfullyExchangedEvents(ch)
		},
		t.cfg.Backoff.NormalPeriod,
		t.cfg.Backoff.MinAbnormalPeriod,
		t.cfg.Backoff.MaxAbnormalPeriod,
	)
}

func (t *MarketPlaceTracker) TrackNftRequestCreatedEvents(ch chan<- etherdata.NFTRequestCreatedEvent) {
	listener := t.listener.WithAddress(t.marketplaceAddress)
	running.WithBackOff(
		t.ctx,
		t.log,
		t.cfg.Prefix+nftRequestCreatedSuffix,
		func(ctx context.Context) error {
			startBlock := uint64(t.network.FirstBlock)
			block, err := t.database.Blocks().FilterByContractAddress(t.marketplaceAddress.String()).FilterByChainId(t.network.ChainId).Get()
			if err != nil {
				return errors.Wrap(err, "failed to get block to begin with")
			}
			if block != nil {
				startBlock = block.NftRequestCreatedBlock
			}

			t.log.WithFields(logan.F{
				"marketplace_address": t.marketplaceAddress,
				"network":             t.network.Name,
			}).Info("start tracking mint events from block ", startBlock)
			return listener.
				From(startBlock).
				WithCtx(ctx).
				WatchNftRequestCreatedEvents(ch)
		},
		t.cfg.Backoff.NormalPeriod,
		t.cfg.Backoff.MinAbnormalPeriod,
		t.cfg.Backoff.MaxAbnormalPeriod,
	)
}

func (t *MarketPlaceTracker) TrackNftRequestCanceledEvents(ch chan<- etherdata.NFTRequestCanceledEvent) {
	listener := t.listener.WithAddress(t.marketplaceAddress)
	running.WithBackOff(
		t.ctx,
		t.log,
		t.cfg.Prefix+nftRequestCanceledSuffix,
		func(ctx context.Context) error {
			startBlock := uint64(t.network.FirstBlock)
			block, err := t.database.Blocks().FilterByContractAddress(t.marketplaceAddress.String()).FilterByChainId(t.network.ChainId).Get()
			if err != nil {
				return errors.Wrap(err, "failed to get block to begin with")
			}
			if block != nil {
				startBlock = block.NftRequestCanceledBlock
			}

			t.log.WithFields(logan.F{
				"marketplace_address": t.marketplaceAddress,
				"network":             t.network.Name,
			}).Info("start tracking mint events from block ", startBlock)
			return listener.
				From(startBlock).
				WithCtx(ctx).
				WatchNftRequestCanceledEvents(ch)
		},
		t.cfg.Backoff.NormalPeriod,
		t.cfg.Backoff.MinAbnormalPeriod,
		t.cfg.Backoff.MaxAbnormalPeriod,
	)
}
