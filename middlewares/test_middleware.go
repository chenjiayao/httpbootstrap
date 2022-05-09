package middlewares

import (
	"httpbootstrap/utils/logger"

	"github.com/gin-gonic/gin"
)

func TestMiddleware(c *gin.Context) {
	//中间件
	logger.Logger.Info("test middleware")
	c.Next()
}
