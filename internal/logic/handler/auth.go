package handler

import (
	"net/http"

	"github.com/sinomoe/fiber/internal/logic/service"

	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/internal/logic/dto"
)

func Auth(ctx *gin.Context) {
	var (
		req  dto.AuthRequest
		resp dto.AuthResponse
		err  error
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if resp.Token, err = service.GetAuth().BuildToken(req.Username); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func ParseToken(ctx *gin.Context) {
	var (
		req  dto.ParseTokenRequest
		resp dto.ParseTokenResponse
		err  error
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if resp, err = service.GetAuth().ParseToken(req.Token); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
