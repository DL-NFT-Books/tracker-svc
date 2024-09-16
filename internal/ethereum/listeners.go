package ethereum

import (
	"context"
	"time"

	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"github.com/ethereum/go-ethereum/common"
)

type (
	TokenListener interface {
		From(from uint64) TokenListener
		To(to uint64) TokenListener
		WithMaxDepth(maxDepth uint64) TokenListener
		WithAddress(address common.Address) TokenListener
		WithCtx(ctx context.Context) TokenListener
		WithDelayBetweenIntervals(delay time.Duration) TokenListener

		WatchTokenSuccessfullyPurchasedEvents(ch chan<- etherdata.TokenSuccessfullyPurchasedEvent) error
		WatchTokenSuccessfullyExchangedEvents(ch chan<- etherdata.TokenSuccessfullyExchangedEvent) error
		WatchNftRequestCreatedEvents(ch chan<- etherdata.NFTRequestCreatedEvent) error
		WatchNftRequestCanceledEvents(ch chan<- etherdata.NFTRequestCanceledEvent) error
	}
)
