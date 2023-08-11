package http

import (
	"github.com/Hank-Kuo/go-kafka-example/config"
	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type Middleware struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewares(cfg *config.Config, logger logger.Logger) *Middleware {
	return &Middleware{cfg: cfg, logger: logger}
}

func NewGlobalMiddlewares(engine *gin.Engine) {
	engine.Use(corsMiddleware())
	engine.NoMethod(httpNotFound)
	engine.NoRoute(httpNotFound)
	Healthz(engine)

	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(engine)

}
