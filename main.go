//go:build go1.8
// +build go1.8

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiru997/go-ex/configs"
	"github.com/kiru997/go-ex/internal"
	"github.com/kiru997/go-ex/monitoring"
	"github.com/kiru997/go-ex/pkg/logger"
	"github.com/kiru997/go-ex/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	cfg, err := configs.NewConfig("./app.yml")
	if err != nil {
		log.Panic("err load config:", err)
	}
	zapLogger := logger.NewZapLogger(cfg.Logger.LogLevel, cfg.Env == "local")
	opts := []fx.Option{
		fx.Provide(func(cfg *configs.BaseConfig) *zap.Logger {
			return zapLogger.With(zap.String("serviceName", cfg.Name))
		}),
		fx.Provide(
			func() *configs.Config {
				return cfg
			},
			func(cfg *configs.Config) *configs.BaseConfig {
				return &cfg.BaseConfig
			},
		),
		fx.Provide(func(r *gin.Engine) *gin.RouterGroup {
			g := r.Group("")
			return g
		}),
		monitoring.ModuleMonitoring,
		internal.Module,
		fx.Provide(server.InitGinEngine),
		fx.Invoke(server.RunHTTPServer),
	}

	err = fx.ValidateApp(opts...)
	if err != nil {
		log.Panic("err provide autowire", err)
	}

	fx.New(
		opts...,
	).Run()
}
