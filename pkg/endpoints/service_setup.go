package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetupRequest struct {
	Password string `form:"password"`
	Phone    string `form:"phone"`
	Code     string `form:"code"`
}

func (s *Service) Setup(ctx *gin.Context) {
	var req SetupRequest

	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	auth := s.tg.Auth()

	if req.Phone != "" {
		auth.SetPhone(req.Phone)
	}

	if req.Code != "" {
		auth.SetCode(req.Code)
	}

	if req.Password != "" {
		auth.SetPassword(req.Password)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"state": auth.StateType(),
	})
}
