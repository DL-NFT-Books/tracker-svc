package runners

import (
	"context"
	"gitlab.com/tokend/nft-books/contract-tracker/internal/config"
)

type Runner func(ctx context.Context)

func initializeRunners(cfg config.Config) (runners []Runner) {
	runners = append(runners, NewFactoryTracker(cfg).Run)

	return
}

func Run(cfg config.Config, ctx context.Context) {
	for _, runner := range initializeRunners(cfg) {
		runner(ctx)
	}
}