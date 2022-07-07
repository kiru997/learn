package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kiru997/go-ex/configs"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InitGinEngine(
	cfg *configs.BaseConfig,
	zapLogger *zap.Logger,
) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if !cfg.Debug {
		r.Use(ginzap.RecoveryWithZap(zapLogger, true))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	r.Use(otelgin.Middleware(cfg.Name))
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "DeviceID", "Accept-Language"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))
	r.Use(requestid.New())

	return r
}

func RunHTTPServer(
	lifecycle fx.Lifecycle,
	zapLogger *zap.Logger,
	cfg *configs.BaseConfig,
	r *gin.Engine,
	p *profiler.Profiler,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					err := r.Run(cfg.Port)
					if err != nil {
						zapLogger.Fatal("r.Run", zap.String("port", cfg.Port), zap.Error(err))
					}
				}()

				zapLogger.Info("HTTP server is running", zap.String("port", cfg.Port), zap.String("env", cfg.Env), zap.Bool("debug", cfg.Debug))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				zapLogger.Info("shutting down server...")

				if p != nil {
					err := p.Stop()
					if err != nil {
						zapLogger.Error("err stop profiler", zap.Error(err))
					}
				}

				return nil
			},
		},
	)
}
