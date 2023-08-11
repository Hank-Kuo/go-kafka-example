package grpc_response

import (
	"github.com/Hank-Kuo/go-kafka-example/pkg/customError/grpcError"
	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"

	"google.golang.org/grpc/status"
)

func Fail(err error, logger logger.Logger) error {
	parseErr := grpcError.ParseError(err)
	if parseErr.Detail != nil {
		logger.Error(parseErr.Detail)
	}
	return status.Error(parseErr.Code, parseErr.GetMessage())
}
