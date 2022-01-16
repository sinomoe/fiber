package helper

import (
	"github.com/gin-gonic/gin"
)

func SetUser(ctx *gin.Context, username string) {
	ctx.Set("username", username)
}

func MustGetUser(ctx *gin.Context) string {
	user := ctx.GetString("username")
	if len(user) == 0 {
		panic("username is empty")
	}
	return user
}
