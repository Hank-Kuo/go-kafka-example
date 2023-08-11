package product

import (
	"context"

	"github.com/Hank-Kuo/go-kafka-example/internal/models"
	"github.com/Hank-Kuo/go-kafka-example/pkg/tracer"
	"github.com/pkg/errors"
)

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	ctx, span := tracer.NewSpan(ctx, "UserRepo.Create", nil)
	defer span.End()

	sqlQuery := `INSERT INTO users(email, password, name) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, sqlQuery, user.Email, user.Password, user.Name)

	if err != nil {
		return errors.Wrap(err, "UserRepo.Create")
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {

	sqlQuery := `
	UPDATE users 
		SET name = COALESCE(NULLIF($2, ''), name),
			password = COALESCE(NULLIF($3, ''), password),
			status = $4
		WHERE email = $1`

	if _, err := r.db.ExecContext(ctx, sqlQuery, user.Email, user.Name, user.Email, user.Status); err != nil {
		return errors.Wrap(err, "userRepo.Update")
	}

	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "UserRepo.GetByEmail", nil)
	defer span.End()

	user := models.User{}
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserRepo.GetByEmail")
	}

	return &user, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "UserRepo.GetAll", nil)
	defer span.End()

	users := []*models.User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "UserRepo.GetAll")
	}

	return users, nil
}
