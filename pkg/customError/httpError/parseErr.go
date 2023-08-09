package httpError

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-kafka-example/pkg/customError"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

var DBSTATE = "pq: "
var SQLSTATE = "sql: "

func unwrap(err error) error {
	for {
		if err == nil {
			return nil
		}
		unwrappedErr := errors.Unwrap(err)
		if unwrappedErr == nil {
			return err
		}
		err = unwrappedErr
	}
}

func ParseError(err error) *Err {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return NewRequestTimeoutError(err)

	case errors.Is(err, customError.ErrOTPCodeNotMatched):
		return &Err{
			Status: http.StatusUnauthorized, Message: customError.ErrOTPCodeNotMatched.Error(),
		}
	case errors.Is(err, customError.ErrEmailNotActive):
		return &Err{
			Status: http.StatusForbidden, Message: customError.ErrEmailNotActive.Error(),
		}
	case errors.Is(err, customError.ErrPasswordCodeNotMatched):
		return &Err{
			Status: http.StatusUnauthorized, Message: customError.ErrPasswordCodeNotMatched.Error(),
		}
	case errors.Is(err, customError.ErrEmailAlreadyActived):
		return &Err{
			Status: http.StatusConflict, Message: customError.ErrEmailAlreadyActived.Error(),
		}
	case strings.Contains(err.Error(), "Field validation"):
		return parseValidatorError(err)
	case strings.Contains(err.Error(), "unmarshal"):
		return parseMarshalError(err)
	case strings.Contains(err.Error(), SQLSTATE):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), DBSTATE):
		return parseDBErrors(err)
	default:
		if restErr, ok := err.(*Err); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseValidatorError(err error) *Err {
	var ve validator.ValidationErrors
	errMessage := "The Field: "
	if errors.As(err, &ve) {
		for i, fe := range ve {
			if i == len(ve)-1 {
				errMessage += fmt.Sprint(strings.ToLower(fe.Field()), " is required")
			} else {
				errMessage += fmt.Sprint(strings.ToLower(fe.Field()), ", ")
			}
		}
	}

	return NewError(http.StatusBadRequest, errMessage, err)
}

func parseMarshalError(err error) *Err {
	errMessage := "The Field: "
	if e, ok := err.(*json.UnmarshalTypeError); ok {
		errMessage += fmt.Sprint(e.Field, " should be ", e.Type)
		return NewError(http.StatusBadRequest, errMessage, err)
	}

	return NewBadRequestError(err)
}

func parseDBErrors(err error) *Err {
	unWrapErr := unwrap(err)
	if pqErr, ok := unWrapErr.(*pq.Error); ok {
		// duplicated
		if pqErr.Code == "23505" {
			errMessage := fmt.Sprint("This Field ", pqErr.Constraint, " is used")
			return NewError(http.StatusConflict, errMessage, err)
		}

	}
	return NewInternalServerError(err)
}

func parseSqlErrors(err error) *Err {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewNotFoundError(err)
	case errors.Is(err, sql.ErrConnDone):
		return NewError(http.StatusInternalServerError, "db connection is done", err)
	case errors.Is(err, sql.ErrNoRows):
		return NewNotFoundError(err)
	case errors.Is(err, sql.ErrTxDone):
		return NewError(http.StatusInternalServerError, "transaction operation error", err)
	default:
		return NewInternalServerError(err)
	}

}
