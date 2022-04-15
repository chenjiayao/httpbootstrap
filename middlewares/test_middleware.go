package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func TestMiddleware(c *gin.Context) {
	//中间件
	log.Info().Msg("test middleware")
	c.Next()
}
