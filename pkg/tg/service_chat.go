package tg

import (
	"context"

	"github.com/opoccomaxao/tgx-api/pkg/models"
	"github.com/pkg/errors"
	"github.com/zelenin/go-tdlib/client"
)

type ChatRequest struct {
	ChatID   int64
	Username string
}

//nolint:cyclop
func (s *Service) GetLinkedChatID(
	_ context.Context,
	req ChatRequest,
) (int64, error) {
	if req.ChatID == 0 && req.Username == "" {
		return 0, errors.WithStack(models.ErrInvalidParams)
	}

	instance, err := s.Client()
	if err != nil {
		return 0, err
	}

	var chat *client.Chat

	switch {
	case req.ChatID != 0:
		chat, err = instance.GetChat(&client.GetChatRequest{
			ChatId: req.ChatID,
		})
	case req.Username != "":
		chat, err = instance.SearchPublicChat(&client.SearchPublicChatRequest{
			Username: req.Username,
		})
	default:
		err = models.ErrNotFound
	}

	if err != nil {
		return 0, errors.WithStack(err)
	}

	if chat == nil {
		return 0, errors.WithStack(models.ErrNotFound)
	}

	typ := chat.Type.ChatTypeType()
	if typ != client.TypeChatTypeSupergroup {
		return 0, errors.Wrapf(models.ErrInvalidParams, "chat type %s", typ)
	}

	sgType, ok := chat.Type.(*client.ChatTypeSupergroup)
	if !ok {
		return 0, errors.WithStack(models.ErrNotFound)
	}

	sg, err := instance.GetSupergroupFullInfo(&client.GetSupergroupFullInfoRequest{
		SupergroupId: sgType.SupergroupId,
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return sg.LinkedChatId, nil
}
