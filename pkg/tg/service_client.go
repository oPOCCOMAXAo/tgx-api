package tg

import (
	"github.com/opoccomaxao/tgx-api/pkg/models"
	"github.com/pkg/errors"
	"github.com/zelenin/go-tdlib/client"
)

func (s *Service) Client() (*client.Client, error) {
	res := s.client
	if res == nil {
		return nil, errors.WithStack(models.ErrNotReady)
	}

	return res, nil
}

func (s *Service) WithClient(visitor func(*client.Client) error) error {
	c, err := s.Client()
	if err != nil {
		return err
	}

	return visitor(c)
}
