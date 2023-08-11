package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"

	"github.com/Hank-Kuo/go-kafka-example/config"
	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"
)

type Server struct {
	engine      *gin.Engine
	grpcEngine  *grpc.Server
	cfg         *config.Config
	db          *sqlx.DB
	kakfaWriter *kafka.Writer
	logger      logger.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, kakfaWriter *kafka.Writer, logger logger.Logger) *Server {
	return &Server{engine: nil, grpcEngine: nil, cfg: cfg, db: db, kakfaWriter: kakfaWriter, logger: logger}
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	httpServer := s.newHttpServer()

	grpcClose, grpcServer, err := s.newGrpcServer()
	if err != nil {
		s.logger.Fatalf("Error gprc serve: %s", err)
	}
	defer grpcClose()

	// graceful shutdown
	<-ctx.Done()
	s.logger.Info("Shutdown Server ...")

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}
