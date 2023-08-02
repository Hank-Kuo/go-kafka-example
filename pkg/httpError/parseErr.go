package httpError

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

func ParseError(err error) *Err {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewNotFoundError(err)
	case strings.Contains(err.Error(), "duplicate"):
		return NewError(http.StatusBadRequest, "Duplicate email", err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRequestTimeoutError(err)
	case strings.Contains(err.Error(), "Field validation"):
		return parseValidatorError(err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewBadRequestError(err)
	case strings.Contains(err.Error(), "UUID"):
		return NewError(http.StatusBadRequest, err.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewUnauthorizedError(err)
	default:
		if restErr, ok := err.(*Err); ok {
			return restErr
		}

		return NewInternalServerError(err)
	}
}

func parseValidatorError(err error) *Err {
	if strings.Contains(err.Error(), "Password") {
		return NewError(http.StatusBadRequest, "Invalid password", err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewError(http.StatusBadRequest, "Invalid email", err)
	}
	return NewBadRequestError(err)
}

func parseSqlErrors(err error) *Err {
	if strings.Contains(err.Error(), "23505") {
		return NewError(http.StatusBadRequest, ExistsEmailError.Error(), err)
	}
	return NewBadRequestError(err)
}
