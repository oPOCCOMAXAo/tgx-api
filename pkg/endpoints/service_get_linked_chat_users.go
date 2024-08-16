package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/tgx-api/pkg/tg"
)

type GetLinkedChatUsersRequest struct {
	Username string `json:"username"`
	ChatID   int64  `json:"chat_id"`
}

type GetLinkedChatUsersResponse struct {
	UserIDs []int64 `json:"user_ids"`
}

func (s *Service) GetLinkedChatUsers(ctx *gin.Context) {
	var req GetLinkedChatUsersRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	chatID, err := s.tg.GetLinkedChatID(ctx, tg.ChatRequest{
		Username: req.Username,
		ChatID:   req.ChatID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	users, err := s.tg.GetChatMembersID(ctx, chatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, &GetLinkedChatUsersResponse{
		UserIDs: users,
	})
}
