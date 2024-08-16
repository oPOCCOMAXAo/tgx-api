package app

import (
	"context"

	"github.com/opoccomaxao/tgx-api/pkg/config"
	"github.com/opoccomaxao/tgx-api/pkg/endpoints"
	"github.com/opoccomaxao/tgx-api/pkg/server"
	"github.com/opoccomaxao/tgx-api/pkg/tg"
	"github.com/samber/do"
)

func InitDependencies(
	cfg config.Config,
) *do.Injector {
	injector := do.New()

	tg.Provide(injector, cfg.TG)
	server.Provide(injector, cfg.Server)
	endpoints.Provide(injector)

	return injector
}

func LaunchDependencies(
	ctx context.Context,
	injector *do.Injector,
	cancelCauseFn func(error),
) error {
	tg, err := tg.Invoke(injector)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	server, err := server.Invoke(injector)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	_, err = endpoints.Invoke(injector)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	for _, serveFn := range []func(context.Context) error{
		tg.Serve,
		server.Serve,
	} {
		go func(serveFn func(context.Context) error) {
			err := serveFn(ctx)
			if err != nil {
				cancelCauseFn(err)
			}
		}(serveFn)
	}

	return nil
}
