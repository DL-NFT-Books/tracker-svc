package consumers

import (
	"context"
	"github.com/dl-nft-books/book-svc/connector/models"
	coreResources "github.com/dl-nft-books/core-svc/resources"
	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"strconv"
)

func (c *MarketPlaceConsumer) ConsumeNFTRequestCreatedEvent(ch <-chan etherdata.NFTRequestCreatedEvent) {
	for {
		select {
		case event := <-ch:
			running.UntilSuccess(
				c.ctx,
				c.logger,
				c.cfg.Prefix+nftRequestCreatedConsumerSuffix,
				func(ctx context.Context) (ok bool, err error) {
					logField := logan.F{"contract_address": c.network.FactoryAddress}

					book, err := c.booker.ListBooks(models.ListBooksParams{
						Contract: []string{event.TokenContract.String()},
						ChainId:  []int64{c.network.ChainId},
					})
					if err != nil {
						c.logger.Error("failed get book")
						return false, errors.Wrap(err, "failed get book", logField)
					}
					if len(book.Data) == 0 {
						c.logger.Warn("could not find book")
						return true, nil
					}

					bookId, err := strconv.Atoi(book.Data[0].ID)
					if err != nil {
						c.logger.Error("failed to parse book id")
						return false, errors.Wrap(err, "failed to parse book id", logField)
					}
					err = c.core.CreateNFTRequest(coreResources.CreateNftRequestAttributes{
						Requester:            event.Requester.String(),
						NftId:                event.NftId.Int64(),
						NftAddress:           event.NftContract.String(),
						ChainId:              c.network.ChainId,
						MarketplaceRequestId: event.RequestId.Int64(),
						BookId:               int64(bookId),
					})
					if err != nil {
						return false, errors.Wrap(err, "failed to consume nft request created transaction", logField)
					}

					// Updating contract`s last nft request created block
					if err = c.database.Blocks().UpdateNFTRequestCreatedBlockColumn(event.BlockNumber, c.network.ChainId); err != nil {
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
