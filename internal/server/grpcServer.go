package server

import (
	"net"
	"time"

	userDelivery "go-kafka-example/internal/api/delivery/user"
	userRepository "go-kafka-example/internal/api/repository/user"
	userService "go-kafka-example/internal/api/service/user"
	userPb "go-kafka-example/pb/user"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *Server) newGrpcServer() (func() error, *grpc.Server, error) {
	l, err := net.Listen("tcp", ":"+s.cfg.Server.GrpcPort)
	if err != nil {
		return nil, nil, errors.Wrap(err, "grpc.net.Listen")
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	grpc_prometheus.Register(grpcServer)

	userRepo := userRepository.NewRepo(s.db, s.kakfaWriter)
	userSrv := userService.NewService(s.cfg, userRepo, s.kakfaWriter, s.logger)
	userHandler := userDelivery.NewGrpcHandler(userSrv, s.logger)
	userPb.RegisterUserServiceServer(grpcServer, userHandler)

	if s.cfg.Server.Debug {
		reflection.Register(grpcServer)
	}

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			s.logger.Fatalf("Error gprc serve: %s", err)
		}
	}()

	return l.Close, grpcServer, nil
}
