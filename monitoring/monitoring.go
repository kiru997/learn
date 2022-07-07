package monitoring

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kiru997/go-ex/configs"
	"github.com/kiru997/go-ex/pkg/tracingutil"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.uber.org/fx"
)

var ModuleMonitoring = fx.Options(
	fx.Provide(
		tracingutil.InitTelemetry,
		NewProfiler,
	),
	fx.Invoke(
		RegisterMetricsExporter,
	),
)

func RegisterMetricsExporter(r *gin.Engine, pe *prometheus.Exporter) {
	if pe != nil {
		r.GET("/metrics", func(c *gin.Context) {
			pe.ServeHTTP(c.Writer, c.Request)
		})
	}
}

func NewProfiler(cfg *configs.BaseConfig) (*profiler.Profiler, error) {
	if !cfg.RemoteProfiler.Enabled {
		return nil, nil
	}

	return profiler.Start(profiler.Config{
		ApplicationName: fmt.Sprintf("learn.%s", cfg.Name),
		ServerAddress:   cfg.RemoteProfiler.ProfilerURL,
	})
}
