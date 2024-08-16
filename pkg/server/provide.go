package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func Provide(
	i *do.Injector,
	cfg Config,
) {
	do.ProvideNamed(i, "server", func(*do.Injector) (*Server, error) {
		res := New(cfg)

		return res, nil
	})
}

func Invoke(i *do.Injector) (*Server, error) {
	return do.InvokeNamed[*Server](i, "server")
}

func InvokeRouter(i *do.Injector) (gin.IRouter, error) {
	server, err := Invoke(i)
	if err != nil {
		return nil, err
	}

	return server.engine, nil
}
