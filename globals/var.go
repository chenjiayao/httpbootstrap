package globals

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	Engine      *gin.Engine   = nil
	RedisClient *redis.Client = nil
	HttpServer  *http.Server  = nil
)
