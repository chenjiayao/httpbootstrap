package globals

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB           = nil
	Engine      *gin.Engine        = nil
	RedisClient *redis.Client      = nil
	HttpServer  *http.Server       = nil
	Logger      *zap.Logger        = nil
	SugarLogger *zap.SugaredLogger = nil
)
