package product

import (
	"context"
	"fmt"

	"github.com/Hank-Kuo/go-kafka-example/pkg/tracer"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (r *userRepo) PublishEmail(ctx context.Context, name, email string) error {
	ctx, span := tracer.NewSpan(ctx, "UserRepo.PublishEmail", nil)
	defer span.End()

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, email, name)),
	}

	err := r.kakfaWriter.WriteMessages(ctx, msg)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "UserRepo.PublishEmail")
	}
	return nil
}
