package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/pkg/base"
	"github.com/sinomoe/fiber/pkg/queue"
	"log"
	"net/http"
)

type SendRequest struct {
	To      string `json:"to" binding:"required"`
	From    string `json:"from" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func main() {
	q := queue.NewRedis("127.0.0.1:6379", "", "mystream", "g1", 0)
	defer q.Shutdown()
	r := gin.Default()
	r.POST("/send", func(c *gin.Context) {
		var (
			req SendRequest
			err error
		)
		if err = c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err = q.Produce(base.Message{
			From:    req.From,
			To:      req.To,
			Message: req.Message,
		}); err != nil {
			log.Println("produce failed err:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, nil)
	})
	r.Run(":8859")
}
