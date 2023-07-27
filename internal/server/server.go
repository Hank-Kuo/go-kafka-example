package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"

	"go-kafka-example/config"
	"go-kafka-example/internal/middleware"
	"go-kafka-example/pkg/logger"
)

type Server struct {
	engine      *gin.Engine
	cfg         *config.Config
	db          *sqlx.DB
	kakfaWriter *kafka.Writer
	logger      logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, kakfaWriter *kafka.Writer, logger logger.Logger) *Server {
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return &Server{engine: gin.Default(), cfg: cfg, db: db, kakfaWriter: kakfaWriter, logger: logger}
}

func (s *Server) Run() error {
	middleware.NewGlobalMiddlewares(s.engine)

	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Port,
		Handler:        s.engine,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.registerHttpHanders()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error http ListenAndServe: %s", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Info("Shutdown Server ...")

	// handle rest request when server grace shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
