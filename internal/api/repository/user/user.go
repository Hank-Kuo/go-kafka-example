package product

import (
	"context"

	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/tracer"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	GetAll(ctx context.Context) ([]*models.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	return nil
}

func (r *userRepo) Update(ctx context.Context, user models.User) error {
	return nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.GetAll", nil)
	defer span.End()

	users := []*models.User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.GetAll")
	}

	return users, nil
}
