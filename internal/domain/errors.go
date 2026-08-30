package domain

import "net/http"

type AppError struct {
	Status    int
	ErrorCode string
	Message   string
}

func (e *AppError) Error() string {
	return e.Message
}

func newErr(status int, code, message string) *AppError {
	return &AppError{Status: status, ErrorCode: code, Message: message}
}

var (
	ErrInvalidRequest        = newErr(http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
	ErrUnauthorized          = newErr(http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	ErrForbidden             = newErr(http.StatusForbidden, "FORBIDDEN", "forbidden")
	ErrNotFound              = newErr(http.StatusNotFound, "NOT_FOUND", "not found")
	ErrConflict              = newErr(http.StatusConflict, "CONFLICT", "conflict")
	ErrDeviceLimitExceeded   = newErr(http.StatusConflict, "DEVICE_LIMIT_EXCEEDED", "device limit exceeded")
	ErrGuardianLimitExceeded = newErr(http.StatusConflict, "GUARDIAN_LIMIT_EXCEEDED", "guardian limit exceeded")
	ErrInvalidVerifyCode     = newErr(http.StatusBadRequest, "INVALID_VERIFY_CODE", "invalid or expired verification code")
	ErrEmailNotVerified      = newErr(http.StatusForbidden, "EMAIL_NOT_VERIFIED", "email not verified")
	ErrEmailAlreadyVerified  = newErr(http.StatusConflict, "EMAIL_ALREADY_VERIFIED", "email already verified")
	ErrInternalError         = newErr(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
)
