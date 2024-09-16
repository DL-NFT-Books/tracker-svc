package token_listeners

import (
	"time"

	"github.com/dl-nft-books/tracker-svc/internal/data/etherdata"
	"github.com/dl-nft-books/tracker-svc/internal/ethereum"
	"github.com/dl-nft-books/tracker-svc/internal/helpers"
	"github.com/dl-nft-books/tracker-svc/solidity/generated/marketplace"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (l *tokenListener) readTokenSuccessfullyExchangedInterval(interval helpers.Interval, ch chan<- etherdata.TokenSuccessfullyExchangedEvent) error {
	instance, err := l.getRPCInstance(*l.address)
	if err != nil {
		return errors.Wrap(err, "failed to get instance")
	}

	iterator, err := instance.FilterTokenSuccessfullyExchanged(
		&bind.FilterOpts{
			Start: interval.From,
			End:   &interval.To,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to initialize an iterator")
	}
	if iterator == nil {
		return errors.From(ethereum.NullIteratorErr, logan.F{
			"contract": l.address.String(),
			"depth":    interval.To - interval.From,
		})
	}

	defer func(iterator *marketplace.MarketplaceTokenSuccessfullyExchangedIterator) {
		if tempErr := iterator.Close(); tempErr != nil {
			err = tempErr
		}
	}(iterator)

	for iterator.Next() {
		raw := iterator.Event
		if raw != nil {
			var event *etherdata.TokenSuccessfullyExchangedEvent
			event, err = l.converter.TokenSuccessfullyExchanged(*raw)
			if err != nil {
				return errors.Wrap(err, "failed to convert raw event to the needed format")
			}

			ch <- *event
		}
	}

	return nil
}

func (l *tokenListener) readTokenSuccessfullyExchangedEvents(ch chan<- etherdata.TokenSuccessfullyExchangedEvent) (err error) {
	lastChainBlock, err := l.rpc.BlockNumber(l.ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get last block in chain")
	}
	if l.maxDepth == nil {
		return l.readTokenSuccessfullyExchangedInterval(helpers.Interval{
			From: *l.from,
			To:   lastChainBlock,
		}, ch)
	}

	if l.to == nil {
		l.to = &lastChainBlock
	}

	for _, interval := range helpers.SplitIntoIntervals(*l.from, *l.to, *l.maxDepth) {
		if err = l.readTokenSuccessfullyExchangedInterval(interval, ch); err != nil {
			return errors.Wrap(err, "failed to read mint interval", logan.F{
				"from": interval.From,
				"to":   interval.To,
			})
		}

		if l.delayBetweenIntervals != nil {
			time.Sleep(*l.delayBetweenIntervals)
		}
	}

	return nil
}

func (l *tokenListener) listenTokenSuccessfullyExchangedEvents(ch chan<- etherdata.TokenSuccessfullyExchangedEvent) (err error) {
	opts := bind.WatchOpts{
		Context: l.ctx,
	}

	filterer, err := l.getWSInstance(*l.address)
	if err != nil {
		return errors.Wrap(err, "failed to initialize a filterer")
	}

	eventsChannel := make(chan *marketplace.MarketplaceTokenSuccessfullyExchanged)
	subscription, err := filterer.WatchTokenSuccessfullyExchanged(&opts, eventsChannel)
	if err != nil {
		return errors.Wrap(err, "failed to watch succeesfully minted events")
	}

	for {
		select {
		case err = <-subscription.Err():
			subscription.Unsubscribe()
			return errors.Wrap(err, "failed to listen to a subscription")
		case raw := <-eventsChannel:
			if raw == nil {
				continue
			}

			var convertedEvent *etherdata.TokenSuccessfullyExchangedEvent
			convertedEvent, err = l.converter.TokenSuccessfullyExchanged(*raw)
			if err != nil {
				return errors.Wrap(err, "failed to convert event to the needed type")
			}

			ch <- *convertedEvent
		}
	}
}

func (l *tokenListener) WatchTokenSuccessfullyExchangedEvents(ch chan<- etherdata.TokenSuccessfullyExchangedEvent) (err error) {
	// We firstly range from l.from to the last chain block and then start a listener
	if err = l.validateParameters(); err != nil {
		return err
	}
	if err = l.readTokenSuccessfullyExchangedEvents(ch); err != nil {
		return errors.Wrap(err, "failed to read previous events")
	}
	if err = l.listenTokenSuccessfullyExchangedEvents(ch); err != nil {
		return errors.Wrap(err, "failed to listen to token purchased events")
	}

	return nil
}
