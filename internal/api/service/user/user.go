package product

import (
	"context"

	"go-kafka-example/config"
	userReop "go-kafka-example/internal/api/repository/user"
	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"
)

type Service interface {
	Create(ctx context.Context, user models.User) error
	GetAll(ctx context.Context) ([]*models.User, error)
}

type userSrv struct {
	cfg      *config.Config
	userRepo userReop.Repository
	logger   logger.Logger
}

func NewService(cfg *config.Config, userRepo userReop.Repository, logger logger.Logger) Service {
	return &userSrv{
		cfg:      cfg,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (r *userSrv) Create(ctx context.Context, user models.User) error {
	return nil
}

func (srv *userSrv) GetAll(ctx context.Context) ([]*models.User, error) {

	c, span := tracer.NewSpan(ctx, "UserService.GetAll", nil)
	defer span.End()

	// ctx, cancel := context.WithTimeout(ctx, srv.cfg.Server.ContextTimeout)
	// defer cancel()
	return srv.userRepo.GetAll(c)

}
