package combiners

import (
	"context"
	"github.com/dl-nft-books/tracker-svc/internal/config"
	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"github.com/dl-nft-books/tracker-svc/internal/service/runners/consumers"
	"github.com/dl-nft-books/tracker-svc/internal/service/runners/trackers"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3"
)

type TokenCombiner struct {
	tracker            *trackers.MarketPlaceTracker
	consumer           *consumers.MarketPlaceConsumer
	logger             *logan.Entry
	ctx                context.Context
	marketplaceAddress common.Address
}

func NewTokenCombiner(cfg config.Config, ctx context.Context, deployedToken trackers.DeployedToken) *TokenCombiner {
	return &TokenCombiner{
		tracker:            trackers.NewMarketPlaceTracker(cfg, ctx, deployedToken),
		consumer:           consumers.NewTokenConsumer(cfg, ctx, deployedToken),
		logger:             cfg.Log(),
		ctx:                ctx,
		marketplaceAddress: deployedToken.Address,
	}
}

func (c *TokenCombiner) ProduceAndConsumeTokenSuccessfullyPurchasedEvent() {
	// Running tracker (producer) and consumer with a combinersChannel joining them
	c.logger.Infof("Initializing mint event consumer and producer for %s", c.marketplaceAddress)
	go func() {
		ch := make(chan etherdata.TokenSuccessfullyPurchasedEvent)
		go c.tracker.TrackTokenSuccessfullyPurchasedEvents(ch)
		go c.consumer.ConsumeTokenSuccessfullyPurchasedEvent(ch)
	}()
}

func (c *TokenCombiner) ProduceAndConsumeTokenSuccessfullyExchangedEvent() {
	// Running tracker (producer) and consumer with a combinersChannel joining them
	c.logger.Infof("Initializing token exchanged event consumer and producer for %s", c.marketplaceAddress)
	go func() {
		ch := make(chan etherdata.TokenSuccessfullyExchangedEvent)
		go c.tracker.TrackTokenSuccessfullyExchangedEvents(ch)
		go c.consumer.ConsumeTokenSuccessfullyExchangedEvent(ch)
	}()
}

func (c *TokenCombiner) ProduceAndConsumeNftRequestCreatedEvent() {
	// Running tracker (producer) and consumer with a combinersChannel joining them
	c.logger.Infof("Initializing mint event consumer and producer for %s", c.marketplaceAddress)
	go func() {
		ch := make(chan etherdata.NFTRequestCreatedEvent)
		go c.tracker.TrackNftRequestCreatedEvents(ch)
		go c.consumer.ConsumeNFTRequestCreatedEvent(ch)
	}()
}

func (c *TokenCombiner) ProduceAndConsumeNftRequestCanceledEvent() {
	// Running tracker (producer) and consumer with a combinersChannel joining them
	c.logger.Infof("Initializing mint event consumer and producer for %s", c.marketplaceAddress)
	go func() {
		ch := make(chan etherdata.NFTRequestCanceledEvent)
		go c.tracker.TrackNftRequestCanceledEvents(ch)
		go c.consumer.ConsumeNFTRequestCanceledEvent(ch)
	}()
}

func (c *TokenCombiner) ProduceAndConsumeAllEvents() {
	c.logger.Infof("Initializing all possible consumers and producers for %s", c.marketplaceAddress)
	c.ProduceAndConsumeTokenSuccessfullyPurchasedEvent()
	c.ProduceAndConsumeTokenSuccessfullyExchangedEvent()
	c.ProduceAndConsumeNftRequestCreatedEvent()
	c.ProduceAndConsumeNftRequestCanceledEvent()
}
