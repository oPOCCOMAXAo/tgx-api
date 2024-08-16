package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/opoccomaxao/tgx-api/pkg/server"
	"github.com/opoccomaxao/tgx-api/pkg/tg"
	"github.com/pkg/errors"
)

type Config struct {
	TG     tg.Config     `envPrefix:"TG_"`
	Server server.Config `envPrefix:"SERVER_"`
}

func Load() (Config, error) {
	var cfg Config

	err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef:       false,
		UseFieldNameByDefault: false,
	})
	if err != nil {
		return cfg, errors.WithStack(err)
	}

	return cfg, nil
}
