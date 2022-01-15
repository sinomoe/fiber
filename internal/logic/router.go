package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/internal/logic/handler"
)

func register(r *gin.Engine) {
	r.POST("/send", handler.SendMessage)
}
