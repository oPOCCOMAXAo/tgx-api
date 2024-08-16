package tg

import (
	"context"

	"github.com/opoccomaxao/tgx-api/pkg/models"
	"github.com/pkg/errors"
	"github.com/zelenin/go-tdlib/client"
)

func (s *Service) GetChatMembersID(
	_ context.Context,
	chatID int64,
) ([]int64, error) {
	instance, err := s.Client()
	if err != nil {
		return nil, err
	}

	const limit = 200

	sg, err := instance.GetChat(&client.GetChatRequest{
		ChatId: chatID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sgType, ok := sg.Type.(*client.ChatTypeSupergroup)
	if !ok {
		return nil, errors.WithStack(models.ErrNotFound)
	}

	req := client.GetSupergroupMembersRequest{
		SupergroupId: sgType.SupergroupId,
		Limit:        limit,
	}
	res := make([]int64, 0)

	for {
		members, err := instance.GetSupergroupMembers(&req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, mem := range members.Members {
			user, ok := mem.MemberId.(*client.MessageSenderUser)
			if !ok {
				continue
			}

			res = append(res, user.UserId)
		}

		if len(members.Members) < limit {
			break
		}

		req.Offset += limit
	}

	return res, nil
}
