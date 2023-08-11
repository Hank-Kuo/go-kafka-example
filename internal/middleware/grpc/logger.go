package grpc

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func peerFromContext(ctx context.Context) (string, bool) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", false
	}

	addr, ok := pr.Addr.(*net.TCPAddr)
	if !ok {
		return "", false
	}

	return strings.Split(addr.IP.String(), ":")[0], true
}

func LoggingInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		peer, _ := peerFromContext(ctx)

		if err != nil {
			l.Errorf("ERROR: method: %s, ip: %s, duration: %s", info.FullMethod, peer, duration)
		} else {
			l.Infof("INFO: method: %s, ip: %s, latency: %s", info.FullMethod, peer, duration)
		}

		return resp, err
	}
}
