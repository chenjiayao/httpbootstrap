package main

import (
	"context"
	"fmt"
	"httpbootstrap/globals"
	"httpbootstrap/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func run() {

	initLogger()
	initCache()
	// initDB()
	initHttpEngine()
	routes.LoadRoute()

	go func() {
		port := viper.GetInt("http.port")
		address := viper.GetString("http.address")
		globals.HttpServer = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", address, port),
			Handler: globals.Engine,
		}
		globals.Logger.Sugar().Infof("http server listen on %s:%d", address, port)
		err := globals.HttpServer.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()
	waitSignal()
}

func initLogger() {
	atomicLevel := zap.NewAtomicLevel()
	level := viper.GetString("log.level")
	switch level {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	default:
		log.Fatalf("error: log level %s is not supported", level)
	}
	zap.NewProduction()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), //日志的编码方式，
		zapcore.AddSync(os.Stdout),            //日志的输出位置
		atomicLevel,                           // 日志等级
	)
	globals.Logger = zap.New(zapCore, zap.AddCaller(), zap.Fields(zap.String("appname", viper.GetString("log.appname"))))
	globals.SugarLogger = globals.Logger.Sugar()
}

func initHttpEngine() {
	globals.Engine = gin.Default()

	if !viper.GetBool("degbug") {
		gin.SetMode(gin.ReleaseMode)
	}

	globals.Engine = gin.New()
	globals.Engine.Use(gin.Recovery())

	if viper.GetBool("http.gzip_enable") {
		globals.Engine.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	if viper.GetBool("http.cors_enable") {
		globals.Engine.Use(cors.New(cors.Config{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
}

func initCache() {
	address := viper.GetString("redis.address")
	port := viper.GetInt32("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	globals.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Password: password,
		DB:       db,
	})

	if err := globals.RedisClient.Ping(context.TODO()).Err(); err != nil {
		globals.SugarLogger.DPanic("redis client ping error", zap.Error(err))
	}
}

func initDB() {
	db, err := gorm.Open(mysql.New(mysql.Config{

		DSN:                       viper.GetString("db.dsn"),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		globals.SugarLogger.DPanicf("db connect error", zap.Error(err))
	}
	globals.DB = db
}

func waitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	err := globals.HttpServer.Shutdown(ctx)
	globals.SugarLogger.Infof("http server shutdown: %s", err)
	globals.Logger.Sync()

	if err != nil {
	}
}
