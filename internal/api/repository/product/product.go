package product

import (
	"context"

	"go-kafka-example/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(ctx context.Context, product models.Product) error
	Update(ctx context.Context, product models.Product) error
	GetAll(ctx context.Context) (*[]models.Product, error)
	Delete(ctx context.Context, productId string) error
}

type productRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product models.Product) error {

	// if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
	// 	return nil, errors.Wrap(err, "userRepo.GetByEmail")
	// }
	return nil
}

func (r *productRepo) Update(ctx context.Context, product models.Product) error {

	return nil
}

func (r *productRepo) GetAll(ctx context.Context) (*[]models.Product, error) {
	return nil, nil
}

func (r *productRepo) Delete(ctx context.Context, product string) error {

	return nil
}
