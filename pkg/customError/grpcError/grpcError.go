package grpcError

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type Err struct {
	Code    codes.Code
	Message string
	Detail  interface{}
}

func (e Err) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func (e Err) GetMessage() string {
	if e.Message == "" {
		return "Internal server error"
	}
	return e.Message
}

func (e Err) GetStatus() codes.Code {
	if e.Code == 0 {
		return codes.Internal
	}
	return e.Code
}

func NewError(code codes.Code, message string, err error) *Err {
	return &Err{
		code, message, err,
	}
}
