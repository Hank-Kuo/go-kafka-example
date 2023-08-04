package product

import (
	"context"
)

func (r *userRepo) PublishEmail(ctx context.Context, name, email string) error {
	// c, span := tracer.NewSpan(ctx, "UserRepo.PublishEmail", nil)
	// defer span.End()

	// msg := kafka.Message{
	// 	Value: []byte(fmt.Sprintf(`{"email": "%s", "name": "%s"}`, email, name)),
	// }

	// err := r.kakfaWriter.WriteMessages(c, msg)

	// if err != nil {
	// 	tracer.AddSpanError(span, err)
	// 	return errors.Wrap(err, "UserRepo.PublishEmail")
	// }
	return nil
}
