package app

import (
	"context"
	"os"
	"os/signal"
	"github.com/opoccomaxao/tgx-api/pkg/config"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	ctx, cancelCauseFn := context.WithCancelCause(context.Background())
	defer cancelCauseFn(nil)

	ctx, cancelFn := signal.NotifyContext(ctx, os.Interrupt)
	defer cancelFn()

	injector := InitDependencies(cfg)
	defer injector.Shutdown() //nolint:errcheck

	err = LaunchDependencies(ctx, injector, cancelCauseFn)
	if err != nil {
		return err
	}

	<-ctx.Done()

	//nolint:wrapcheck
	return context.Cause(ctx)
}
