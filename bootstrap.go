package main

import (
	"context"
	"fmt"
	"httptemplate/globals"
	"httptemplate/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

func run() {

	initLog()
	initCache()
	initDB()
	initHttpEngine()
	routes.LoadRoute()

	go func() {
		port := viper.GetInt("http.port")
		address := viper.GetString("http.address")
		globals.HttpServer = &http.Server{
			Addr:    fmt.Sprintf("%s:%d", address, port),
			Handler: globals.Engine,
		}
		log.Info().Msgf("http server start at %s:%d", address, port)

		err := globals.HttpServer.ListenAndServe()
		if err != nil {
			log.Fatal().Msgf("http server start failed: %v", err)
			os.Exit(1)
		}
	}()
	waitSignal()
}

func initLog() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logLevel := viper.GetString("log_level")

	switch logLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

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
		log.Fatal().Msgf("redis connect failed: %v", err)
		os.Exit(1)
	}
}

func initDB() {

}

func waitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	err := globals.HttpServer.Shutdown(ctx)
	log.Info().Msgf("http server shutdown: %v", err)
	if err != nil {
	}
}
