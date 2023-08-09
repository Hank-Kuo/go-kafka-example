package customError

import "errors"

var (
	ErrBadRequest             = errors.New("Bad request")
	ErrWrongCredentials       = errors.New("Wrong Credentials")
	ErrNotFound               = errors.New("Not Found")
	ErrUnauthorized           = errors.New("Unauthorized")
	ErrForbidden              = errors.New("Forbidden")
	ErrPermissionDenied       = errors.New("Permission Denied")
	ErrExpiredCSRFError       = errors.New("Expired CSRF token")
	ErrWrongCSRFToken         = errors.New("Wrong CSRF token")
	ErrCSRFNotPresented       = errors.New("CSRF not presented")
	ErrNotRequiredFields      = errors.New("No such required fields")
	ErrBadQueryParams         = errors.New("Invalid query params")
	ErrInternalServerError    = errors.New("Internal Server Error")
	ErrRequestTimeoutError    = errors.New("Request Timeout")
	ErrExistsEmailError       = errors.New("User with given email already exists")
	ErrInvalidJWTToken        = errors.New("Invalid JWT token")
	ErrInvalidJWTClaims       = errors.New("Invalid JWT claims")
	ErrNotAllowedImageHeader  = errors.New("Not allowed image header")
	ErrContextCancel          = errors.New("Context cacnel")
	ErrDeadlineExceeded       = errors.New("Context exceed")
	ErrOTPCodeNotMatched      = errors.New("OTP code is wrong")
	ErrEmailNotActive         = errors.New("Your email isn't active")
	ErrPasswordCodeNotMatched = errors.New("Password is wrong")
	ErrEmailAlreadyActived    = errors.New("you already activate your email")
)
