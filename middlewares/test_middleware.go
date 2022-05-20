package middlewares

import "github.com/gin-gonic/gin"

func TestMiddleware(c *gin.Context) {

	c.Next()
}
