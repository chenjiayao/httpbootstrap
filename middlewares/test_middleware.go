package middlewares

import (
	"httpbootstrap/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestMiddleware(c *gin.Context) {
	//中间件
	logger.Info("test middleware", zap.String("url", c.Request.URL.String()))

	c.Next()
}
