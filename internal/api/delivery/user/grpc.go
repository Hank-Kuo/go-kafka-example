package user

import (
	"context"

	userSrv "github.com/Hank-Kuo/go-kafka-example/internal/api/service/user"
	"github.com/Hank-Kuo/go-kafka-example/internal/models"
	userPb "github.com/Hank-Kuo/go-kafka-example/pb/user"
	"github.com/Hank-Kuo/go-kafka-example/pkg/customError"
	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"
	"github.com/Hank-Kuo/go-kafka-example/pkg/response/grpc_response"
	"github.com/Hank-Kuo/go-kafka-example/pkg/tracer"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcHandler struct {
	userSrv userSrv.Service
	logger  logger.Logger
}

func NewGrpcHandler(userSrv userSrv.Service, logger logger.Logger) *grpcHandler {
	return &grpcHandler{userSrv: userSrv, logger: logger}
}

func (h *grpcHandler) GetUsers(ctx context.Context, request *userPb.GetUsersRequest) (*userPb.GetUsersResponse, error) {

	ctx, span := tracer.NewSpan(ctx, "UserGrpcHandler.Register", nil)
	defer span.End()

	data, err := h.userSrv.GetAll(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, grpc_response.Fail(err, h.logger)
	}

	user := []*userPb.User{}

	for _, v := range data {
		user = append(user, &userPb.User{
			Id:        v.ID,
			Name:      v.Name,
			Email:     v.Email,
			CreatedAt: timestamppb.New(v.CreatedAt),
		})
	}

	return &userPb.GetUsersResponse{User: user}, nil
}

func (h *grpcHandler) Register(ctx context.Context, request *userPb.RegisterRequest) (*userPb.RegisterResponse, error) {

	ctx, span := tracer.NewSpan(ctx, "UserGrpcHandler.Register", nil)
	defer span.End()

	if err := h.userSrv.Register(ctx, &models.User{
		Name:     request.Name,
		Password: request.Password,
		Email:    request.Email,
	}); err != nil {
		tracer.AddSpanError(span, err)
		return nil, grpc_response.Fail(err, h.logger)
	}

	return &userPb.RegisterResponse{}, nil
}

func (h *grpcHandler) Login(ctx context.Context, request *userPb.LoginRequest) (*userPb.LoginResponse, error) {

	ctx, span := tracer.NewSpan(ctx, "UserGrpcHandler.Login", nil)
	defer span.End()

	user, err := h.userSrv.Login(ctx, request.Email, request.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, grpc_response.Fail(err, h.logger)
	}

	return &userPb.LoginResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (h *grpcHandler) UserActive(ctx context.Context, request *userPb.UserActiveRequest) (*userPb.UserActiveResponse, error) {

	ctx, span := tracer.NewSpan(ctx, "UserGrpcHandler.UserActive", nil)
	defer span.End()

	err := h.userSrv.Active(ctx, request.Email, request.Otp)
	if err != nil {
		return nil, grpc_response.Fail(err, h.logger)
	}

	return nil, nil
}

func (h *grpcHandler) SendEmail(ctx context.Context, request *userPb.SendEmailRequest) (*userPb.SendEmailResponse, error) {

	ctx, span := tracer.NewSpan(ctx, "UserGrpcHandler.SendEmail", nil)
	defer span.End()

	user, err := h.userSrv.GetByEmail(ctx, request.Email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, grpc_response.Fail(err, h.logger)
	}

	if user.Status {
		return nil, grpc_response.Fail(customError.ErrEmailAlreadyActived, h.logger)

	} else {
		if err := h.userSrv.PublishEmail(ctx, user.Name, user.Email); err != nil {
			tracer.AddSpanError(span, err)
			return nil, grpc_response.Fail(err, h.logger)
		}
	}

	return nil, nil
}
