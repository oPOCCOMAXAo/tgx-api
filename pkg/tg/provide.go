package tg

import (
	"github.com/samber/do"
)

func Provide(
	i *do.Injector,
	cfg Config,
) {
	do.ProvideNamed(i, "tg", func(*do.Injector) (*Service, error) {
		res := New(cfg)

		return res, nil
	})
}

func Invoke(i *do.Injector) (*Service, error) {
	return do.InvokeNamed[*Service](i, "tg")
}
