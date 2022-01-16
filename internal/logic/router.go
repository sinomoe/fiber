package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/internal/logic/handler"
	"github.com/sinomoe/fiber/internal/logic/middleware"
)

func register(r *gin.Engine) {
	r.POST("/send", middleware.Auth(), handler.SendMessage)
	r.POST("/auth", handler.Auth)
	r.POST("/auth/parse", handler.ParseToken)
}
