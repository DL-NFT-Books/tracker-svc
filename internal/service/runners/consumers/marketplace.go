package consumers

import (
	"context"
	documenter "github.com/dl-nft-books/blob-svc/connector/api"
	booker "github.com/dl-nft-books/book-svc/connector"
	core "github.com/dl-nft-books/core-svc/connector"
	"github.com/dl-nft-books/network-svc/connector/models"
	"github.com/dl-nft-books/tracker-svc/internal/config"
	"github.com/dl-nft-books/tracker-svc/internal/data"
	"github.com/dl-nft-books/tracker-svc/internal/data/postgres"
	"github.com/dl-nft-books/tracker-svc/internal/helpers"
	"github.com/dl-nft-books/tracker-svc/internal/service/runners/trackers"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	mintConsumerSuffix               = "-marketplace-mint"
	nftRequestCreatedConsumerSuffix  = "-marketplace-nft-req-created"
	nftRequestCanceledConsumerSuffix = "-marketplace-nft-req-canceled"
	nftExchangedConsumerSuffix       = "-marketplace-nft-exchanged"
)

type MarketPlaceConsumer struct {
	logger   *logan.Entry
	cfg      config.Runner
	ctx      context.Context
	database data.DB

	ipfsLoader         *helpers.IpfsLoader
	network            models.NetworkDetailedResponse
	booker             *booker.Connector
	core               *core.Connector
	documenter         *documenter.Connector
	marketplaceAddress common.Address
}

func NewTokenConsumer(cfg config.Config, ctx context.Context, deployedToken trackers.DeployedToken) *MarketPlaceConsumer {
	return &MarketPlaceConsumer{
		logger:     cfg.Log(),
		ctx:        ctx,
		cfg:        cfg.Consumers(),
		database:   postgres.NewDB(cfg.DB()),
		network:    deployedToken.Network,
		ipfsLoader: helpers.NewIpfsLoader(cfg),

		booker:             cfg.BookerConnector(),
		core:               cfg.GeneratorConnector(),
		documenter:         cfg.DocumenterConnector(),
		marketplaceAddress: deployedToken.Address,
	}
}
