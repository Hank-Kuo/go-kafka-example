package product

import (
	"context"
	"fmt"

	"go-kafka-example/config"
	userReop "go-kafka-example/internal/api/repository/user"
	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
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
	if err := srv.userRepo.Create(c, user); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.Register")
	}

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, user.Email, user.Name)),
	}
	err := srv.kakfaWriter.WriteMessages(ctx, msg)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserService.KafkaWriter")
	}
	return nil
}

func (srv *userSrv) PublishEmail(ctx context.Context) []byte {
	_, span := tracer.NewSpan(ctx, "UserService.PublishEmail", nil)
	defer span.End()

	return []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, "", ""))
}
func (srv *userSrv) Login(ctx context.Context, email string, password string) (*models.User, error) {
	c, span := tracer.NewSpan(ctx, "UserService.Login", nil)
	defer span.End()

	user, err := srv.userRepo.GetByEmail(c, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserService.GetByEmail")
	}

	return user, nil
}

func (srv *userSrv) GetAll(ctx context.Context) ([]*models.User, error) {

	c, span := tracer.NewSpan(ctx, "UserService.GetAll", nil)
	defer span.End()

	// ctx, cancel := context.WithTimeout(ctx, srv.cfg.Server.ContextTimeout)
	// defer cancel()
	return srv.userRepo.GetAll(c)

}
