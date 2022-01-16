package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/internal/logic/dto"
	"github.com/sinomoe/fiber/internal/logic/helper"
	"github.com/sinomoe/fiber/internal/logic/service"
)

func SendMessage(c *gin.Context) {
	var (
		req      dto.SendMessageRequest
		resp     dto.SendMessageResponse
		username = helper.MustGetUser(c)
		err      error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if resp, err = service.GetQueue().SendMessage(c, username, req); err != nil {
		log.Println("produce failed err:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, resp)
}
