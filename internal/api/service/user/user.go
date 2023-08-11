package product

import (
	"context"

	"go-kafka-example/config"
	"go-kafka-example/pkg/customError"
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
	Register(ctx context.Context, user *models.User) error
	PublishEmail(ctx context.Context, name, email string) error
	Login(ctx context.Context, eamil, password string) (*dto.LoginResDto, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Active(ctx context.Context, email, optCode string) error
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

func (srv *userSrv) Register(ctx context.Context, user *models.User) error {
	c, span := tracer.NewSpan(ctx, "UserService.Register", nil)
	defer span.End()

	// hashing password before insert into database
	hashPassword, err := utils.HashText(user.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	user.Password = hashPassword

	if err := srv.userRepo.Create(c, *user); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}

	if err := srv.userRepo.PublishEmail(ctx, user.Name, user.Email); err != nil {
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
		return nil, customError.ErrPasswordCodeNotMatched
	}

	if !user.Status {
		return nil, customError.ErrEmailNotActive
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

	if err := srv.userRepo.PublishEmail(ctx, name, email); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	return nil
}

func (srv *userSrv) Active(ctx context.Context, email string, otpCode string) error {
	ctx, span := tracer.NewSpan(ctx, "UserService.PublishEmail", nil)
	defer span.End()

	token := utils.GetToken(email)
	isValid, err := utils.ValidOTP(token, otpCode)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Active")
	}

	if !isValid {
		return customError.ErrOTPCodeNotMatched
	}

	if err := srv.userRepo.Update(ctx, &models.User{Status: true, Email: email}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Active")
	}

	return nil
}

func (srv *userSrv) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "UserService.GetByEmail", nil)
	defer span.End()

	user, err := srv.userRepo.GetByEmail(ctx, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserService.GetByEmail")
	}
	return user, nil
}
