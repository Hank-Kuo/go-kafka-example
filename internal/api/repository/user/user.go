package product

import (
	"context"

	"go-kafka-example/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type Repository interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	PublishEmail(ctx context.Context, name, email string) error
}

type userRepo struct {
	db          *sqlx.DB
	kakfaWriter *kafka.Writer
}

func NewRepo(db *sqlx.DB, kakfaWriter *kafka.Writer) Repository {
	return &userRepo{db: db, kakfaWriter: kakfaWriter}
}
