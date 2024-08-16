package endpoints

import (
	"github.com/opoccomaxao/tgx-api/pkg/server"
	"github.com/opoccomaxao/tgx-api/pkg/tg"
	"github.com/samber/do"
)

//nolint:wrapcheck
func Provide(
	i *do.Injector,
) {
	do.ProvideNamed(i, "endpoints", func(i *do.Injector) (*Service, error) {
		tg, err := tg.Invoke(i)
		if err != nil {
			return nil, err
		}

		res := New(
			tg,
		)

		router, err := server.InvokeRouter(i)
		if err != nil {
			return nil, err
		}

		err = res.Init(router)
		if err != nil {
			return nil, err
		}

		return res, nil
	})
}

func Invoke(i *do.Injector) (*Service, error) {
	return do.InvokeNamed[*Service](i, "endpoints")
}
