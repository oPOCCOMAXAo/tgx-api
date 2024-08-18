package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao/tgx-api/pkg/tg"
)

type Service struct {
	tg *tg.Service
}

func New(
	tg *tg.Service,
) *Service {
	return &Service{
		tg: tg,
	}
}

func (s *Service) Init(
	router gin.IRouter,
) error {
	router.GET("/setup", s.Setup)
	router.POST("/setup", s.Setup)
	router.GET("/linked_chat_users", s.GetLinkedChatUsers)
	router.POST("/linked_chat_users", s.GetLinkedChatUsers)

	return nil
}
