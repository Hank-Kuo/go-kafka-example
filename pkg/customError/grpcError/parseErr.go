package grpcError

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Hank-Kuo/go-kafka-example/pkg/customError"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
)

func unwrap(err error) error {
	for {
		if err == nil {
			return nil
		}
		unwrappedErr := errors.Unwrap(err)
		if unwrappedErr == nil {
			return err
		}
		err = unwrappedErr
	}
}

var DBSTATE = "pq: "
var SQLSTATE = "sql: "

func ParseError(err error) *Err {
	switch {
	case errors.Is(err, context.Canceled):
		return NewError(codes.Canceled, customError.ErrContextCancel.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewError(codes.Canceled, customError.ErrDeadlineExceeded.Error(), err)
	case errors.Is(err, customError.ErrOTPCodeNotMatched):
		return NewError(codes.AlreadyExists, customError.ErrOTPCodeNotMatched.Error(), err)
	case errors.Is(err, customError.ErrEmailNotActive):
		return NewError(codes.PermissionDenied, customError.ErrEmailNotActive.Error(), err)
	case errors.Is(err, customError.ErrPasswordCodeNotMatched):
		return NewError(codes.Unauthenticated, customError.ErrPasswordCodeNotMatched.Error(), err)
	case errors.Is(err, customError.ErrEmailAlreadyActived):
		return NewError(codes.AlreadyExists, customError.ErrEmailAlreadyActived.Error(), err)
	case strings.Contains(err.Error(), SQLSTATE):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), DBSTATE):
		return parseDBErrors(err)
	}
	return NewError(codes.Internal, customError.ErrInternalServerError.Error(), err)

}

func parseDBErrors(err error) *Err {
	unWrapErr := unwrap(err)
	if pqErr, ok := unWrapErr.(*pq.Error); ok {
		// duplicated
		if pqErr.Code == "23505" {
			errMessage := fmt.Sprint("This Field ", pqErr.Constraint, " is used")
			return NewError(codes.AlreadyExists, errMessage, err)
		}
	}
	return NewError(codes.Internal, customError.ErrInternalServerError.Error(), err)
}

func parseSqlErrors(err error) *Err {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewError(codes.NotFound, customError.ErrNotFound.Error(), err)
	case errors.Is(err, sql.ErrConnDone):
		return NewError(http.StatusInternalServerError, "db connection is done", err)
	case errors.Is(err, sql.ErrNoRows):
		return NewError(codes.NotFound, customError.ErrNotFound.Error(), err)
	case errors.Is(err, sql.ErrTxDone):
		return NewError(http.StatusInternalServerError, "transaction operation error", err)
	default:
		return NewError(codes.Internal, customError.ErrInternalServerError.Error(), err)
	}

}
