package tg

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zelenin/go-tdlib/client"
)

type Service struct {
	cfg    Config
	client *client.Client // Deprecated: do not use directly, use Client() instead.
	auth   *authHandler
}

type Config struct {
	APIID   int64  `env:"API_ID,required"`
	APIHash string `env:"API_HASH,required"`
	DataDir string `env:"DATA_DIR,required"`
}

func New(
	cfg Config,
) *Service {
	return &Service{
		cfg:    cfg,
		client: nil,
		auth: newAuthHandler(&client.SetTdlibParametersRequest{
			UseTestDc:           false,
			DatabaseDirectory:   filepath.Join(cfg.DataDir, "database"),
			FilesDirectory:      filepath.Join(cfg.DataDir, "files"),
			UseFileDatabase:     false,
			UseChatInfoDatabase: false,
			UseMessageDatabase:  false,
			UseSecretChats:      false,
			ApiId:               int32(cfg.APIID),
			ApiHash:             cfg.APIHash,
			SystemLanguageCode:  "en",
			DeviceModel:         "server",
			ApplicationVersion:  "1.0",
		}),
	}
}

func (s *Service) Serve(_ context.Context) error {
	var err error

	s.client, err = client.NewClient(s.auth)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Service) Shutdown() error {
	if s.client != nil {
		_, err := s.client.Close()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Service) Auth() AuthHandler {
	return s.auth
}
