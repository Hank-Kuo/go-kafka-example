package product

import (
	"context"

	"go-kafka-example/config"
	productRepository "go-kafka-example/internal/api/repository/product"
	"go-kafka-example/internal/models"
	"go-kafka-example/pkg/logger"
)

type ProductService interface {
	Create(ctx context.Context, product models.Product) error
	GetAll(ctx context.Context) (*[]models.Product, error)
	Delete(ctx context.Context, productId string) error
}

type productSrv struct {
	cfg         *config.Config
	productRepo productRepository.ProductRepository
	logger      logger.Logger
}

func NewService(cfg *config.Config, productRepo productRepository.ProductRepository, logger logger.Logger) ProductService {
	return &productSrv{
		cfg:         cfg,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (r *productSrv) Create(ctx context.Context, product models.Product) error {

	// if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
	// 	return nil, errors.Wrap(err, "userRepo.GetByEmail")
	// }
	return nil
}

func (r *productSrv) GetAll(ctx context.Context) (*[]models.Product, error) {
	return nil, nil
}

func (r *productSrv) Delete(ctx context.Context, product string) error {

	return nil
}
