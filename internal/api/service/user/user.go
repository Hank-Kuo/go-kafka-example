package product

import (
	"context"
	"fmt"

	"go-kafka-example/config"
	"go-kafka-example/pkg/utils"

	userReop "go-kafka-example/internal/api/repository/user"
	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	PublishEmail(ctx context.Context, name, email string) error
	Login(ctx context.Context, eamil, password string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
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

func (srv *userSrv) Register(ctx context.Context, user models.User) error {
	c, span := tracer.NewSpan(ctx, "UserService.Register", nil)
	defer span.End()

	// hashing password before insert into database
	hashPassword, err := utils.HashText(user.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	user.Password = hashPassword

	if err := srv.userRepo.Create(c, user); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}
	return nil
}

func (srv *userSrv) Login(ctx context.Context, email string, password string) (*models.User, error) {
	c, span := tracer.NewSpan(ctx, "UserService.Login", nil)
	defer span.End()

	user, err := srv.userRepo.GetByEmail(c, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserService.Login")
	}

	if err := utils.CheckTextHash(password, user.Password); err != nil {
		if err != nil {
			tracer.AddSpanError(span, err)
			return nil, errors.Wrap(err, "UserService.Login")
		}
	}

	return user, nil
}

func (srv *userSrv) GetAll(ctx context.Context) ([]*models.User, error) {

	c, span := tracer.NewSpan(ctx, "UserService.GetAll", nil)
	defer span.End()

	// ctx, cancel := context.WithTimeout(ctx, srv.cfg.Server.ContextTimeout) defer cancel()
	return srv.userRepo.GetAll(c)

}

func (srv *userSrv) PublishEmail(ctx context.Context, name, email string) error {
	c, span := tracer.NewSpan(ctx, "UserService.PublishEmail", nil)
	defer span.End()

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, email, name)),
	}

	err := srv.kakfaWriter.WriteMessages(c, msg)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.PublishEmail")
	}

	return nil
}
