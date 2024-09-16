package consumers

import (
	"context"
	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

func (c *MarketPlaceConsumer) ConsumeNFTRequestCanceledEvent(ch <-chan etherdata.NFTRequestCanceledEvent) {
	for {
		select {
		case event := <-ch:
			running.UntilSuccess(
				c.ctx,
				c.logger,
				c.cfg.Prefix+nftRequestCanceledConsumerSuffix,
				func(ctx context.Context) (ok bool, err error) {
					logField := logan.F{"contract_address": c.network.FactoryAddress}

					err = c.core.CancelNFTRequest(event.RequestId.Int64())
					if err != nil {
						return false, errors.Wrap(err, "failed to consume nft request created transaction", logField)
					}

					// Updating contract`s last nft request created block
					if err = c.database.Blocks().UpdateNFTRequestCanceledBlockColumn(event.BlockNumber, c.network.ChainId); err != nil {
						return false, errors.Wrap(err, "failed to update contract`s last nft request created block")
					}

					c.logger.WithFields(logField).Infof("Successfully processed nft request created event of a marketplace with id %d", event.RequestId)
					return true, nil
				},
				c.cfg.Backoff.MinAbnormalPeriod,
				c.cfg.Backoff.MaxAbnormalPeriod,
			)
		}
	}
}
