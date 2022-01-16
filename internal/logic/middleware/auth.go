package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/internal/logic/dto"
	"github.com/sinomoe/fiber/internal/logic/helper"
	"github.com/sinomoe/fiber/internal/logic/service"
)

type AuthHeader struct {
	Authentication string `header:"authentication" binding:"required"`
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			req    AuthHeader
			token  string
			claims dto.AuthClaims
			err    error
		)
		if err = ctx.ShouldBindHeader(&req); err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		token = strings.TrimPrefix(req.Authentication, "Bearer ")
		if claims, err = service.GetAuth().ParseToken(token); err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		helper.SetUser(ctx, claims.Username)
	}
}
