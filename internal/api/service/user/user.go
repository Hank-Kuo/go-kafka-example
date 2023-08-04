package product

import (
	"context"
	"fmt"
	"net/http"

	"go-kafka-example/config"
	"go-kafka-example/pkg/httpError"
	"go-kafka-example/pkg/utils"

	userReop "go-kafka-example/internal/api/repository/user"
	"go-kafka-example/internal/dto"
	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	Register(ctx context.Context, user *dto.RegisterReqDto) error
	PublishEmail(ctx context.Context, name, email string) error
	Login(ctx context.Context, eamil, password string) (*dto.LoginResDto, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Active(ctx context.Context, email, optCode string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userSrv struct {
	cfg         *config.Config
	userRepo    userReop.Repository
	kakfaWriter *kafka.Writer
	logger      logger.Logger
}

func NewService(cfg *config.Config, userRepo userReop.Repository, kakfaWriter *kafka.Writer, logger logger.Logger) Service {
	return &userSrv{
		cfg:         cfg,
		userRepo:    userRepo,
		kakfaWriter: kakfaWriter,
		logger:      logger,
	}
}

func (srv *userSrv) Register(ctx context.Context, user *dto.RegisterReqDto) error {
	c, span := tracer.NewSpan(ctx, "UserService.Register", nil)
	defer span.End()

	// hashing password before insert into database
	hashPassword, err := utils.HashText(user.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	user.Password = hashPassword

	if err := srv.userRepo.Create(c, models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	return nil
}

func (srv *userSrv) Login(ctx context.Context, email string, password string) (*dto.LoginResDto, error) {
	ctx, span := tracer.NewSpan(ctx, "UserService.Login", nil)
	defer span.End()

	user, err := srv.userRepo.GetByEmail(ctx, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserService.Login")
	}

	if err := utils.CheckTextHash(password, user.Password); err != nil {
		return nil, errors.New("Field validation: Password error")
	}

	if !user.Status {
		return nil, &httpError.Err{
			Status: http.StatusUnauthorized, Message: "not activate your email", Detail: nil,
		}
	}

	return &dto.LoginResDto{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (srv *userSrv) GetAll(ctx context.Context) ([]*models.User, error) {

	ctx, span := tracer.NewSpan(ctx, "UserService.GetAll", nil)
	defer span.End()

	// ctx, cancel := context.WithTimeout(ctx, srv.cfg.Server.ContextTimeout) defer cancel()
	return srv.userRepo.GetAll(ctx)
}

func (srv *userSrv) PublishEmail(ctx context.Context, name, email string) error {
	ctx, span := tracer.NewSpan(ctx, "UserService.PublishEmail", nil)
	defer span.End()

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, email, name)),
	}

	err := srv.kakfaWriter.WriteMessages(ctx, msg)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.PublishEmail")
	}

	return nil
}

func (srv *userSrv) Active(ctx context.Context, email string, otpCode string) (bool, error) {
	ctx, span := tracer.NewSpan(ctx, "UserService.PublishEmail", nil)
	defer span.End()

	token := utils.GetToken(email)
	isValid, err := utils.ValidOTP(token, otpCode)

	if err != nil {
		tracer.AddSpanError(span, err)
		return false, errors.Wrap(err, "UserService.Active")
	}

	if !isValid {
		return false, &httpError.Err{
			Status: http.StatusBadRequest, Message: "OTP code error", Detail: nil,
		}

	}

	if err := srv.userRepo.Update(ctx, &models.User{Status: true, Email: email}); err != nil {
		tracer.AddSpanError(span, err)
		return false, errors.Wrap(err, "UserService.Active")
	}

	return isValid, nil
}

func (srv *userSrv) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "UserService.GetByEmail", nil)
	defer span.End()
	return srv.userRepo.GetByEmail(ctx, email)
}
