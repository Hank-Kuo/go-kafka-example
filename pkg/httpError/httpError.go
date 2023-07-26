package httpError

import (
	"fmt"
	"net/http"
)

type Err struct {
	Status  int
	Message string
	Detail  interface{}
}

func (e Err) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func (e Err) GetMessage() string {
	if e.Message == "" {
		return "Internal server error"
	}
	return e.Message
}

func (e Err) GetStatus() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func NewError(httpCode int, message string, err error) *Err {
	return &Err{
		httpCode, message, err,
	}
}

func NewNotFoundError(err error) *Err {
	return &Err{
		http.StatusNotFound, NotFound.Error(), err,
	}
}

func NewInternalServerError(err error) *Err {
	return &Err{
		http.StatusInternalServerError, InternalServerError.Error(), err,
	}
}

func NewBadRequestError(err error) *Err {
	return &Err{
		http.StatusBadRequest, BadRequest.Error(), err,
	}
}

func NewUnauthorizedError(err error) *Err {
	return &Err{
		http.StatusUnauthorized, Unauthorized.Error(), err,
	}
}
func NewRequestTimeoutError(err error) *Err {
	return &Err{
		http.StatusRequestTimeout, RequestTimeoutError.Error(), err,
	}
}

func NewBadQueryParamsError(err error) *Err {
	return &Err{
		http.StatusBadRequest, BadQueryParams.Error(), err,
	}
}
