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

	"go-kafka-example/config"
	"go-kafka-example/pkg/logger"
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
	// if cfg.Server.Debug {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// grpcEngine := grpc.NewServer(
	// 	grpc.KeepaliveParams(keepalive.ServerParameters{
	// 		MaxConnectionIdle: maxConnectionIdle * time.Minute,
	// 		Timeout:           gRPCTimeout * time.Second,
	// 		MaxConnectionAge:  maxConnectionAge * time.Minute,
	// 		Time:              gRPCTime * time.Minute,
	// 	}),
	// 	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	// 		grpc_ctxtags.UnaryServerInterceptor(),
	// 		grpc_opentracing.UnaryServerInterceptor(),
	// 		grpc_prometheus.UnaryServerInterceptor,
	// 		grpc_recovery.UnaryServerInterceptor(),
	// 	)),
	// )

	// if cfg.Server.Debug {
	// 	reflection.Register(grpcEngine)
	// }

	return &Server{engine: nil, grpcEngine: nil, cfg: cfg, db: db, kakfaWriter: kakfaWriter, logger: logger}
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// http server
	// middleware.NewGlobalMiddlewares(s.engine)

	// httpServer := &http.Server{
	// 	Addr:           ":" + s.cfg.Server.Port,
	// 	Handler:        s.engine,
	// 	ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
	// 	WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.registerHttpHanders()

	// go func() {
	// 	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		s.logger.Fatalf("Error http ListenAndServe: %s", err)
	// 	}
	// }()
	httpServer := s.newHttpServer()

	// grpc server
	// grpcServer, err := net.Listen("tcp", ":"+s.cfg.Server.GrpcPort)
	// if err != nil {
	// 	s.logger.Fatalf("Error grpc Listen: %s", err)
	// }

	// defer grpcServer.Close()

	// go func() {
	// 	if err := s.grpcEngine.Serve(grpcServer); err != nil {
	// 		s.logger.Fatalf("Error gprc serve: %s", err)
	// 	}
	// }()

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
