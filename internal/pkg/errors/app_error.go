package errors

import "net/http"

type AppError struct {
	Code     Code
	Message  string
	Internal error
	Status   int
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Internal != nil {
		return e.Internal.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Internal
}

func New(code Code, status int, message string, internal error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Internal: internal,
		Status:   status,
	}
}

func BadRequest(message string, internal error) *AppError {
	return New(CodeBadRequest, http.StatusBadRequest, message, internal)
}

func Unauthorized(message string, internal error) *AppError {
	return New(CodeUnauthorized, http.StatusUnauthorized, message, internal)
}

func Forbidden(message string, internal error) *AppError {
	return New(CodeForbidden, http.StatusForbidden, message, internal)
}

func NotFound(message string, internal error) *AppError {
	return New(CodeNotFound, http.StatusNotFound, message, internal)
}

func Conflict(message string, internal error) *AppError {
	return New(CodeConflict, http.StatusConflict, message, internal)
}

func Timeout(message string, internal error) *AppError {
	return New(CodeTimeout, http.StatusGatewayTimeout, message, internal)
}

func ServiceNotReady(message string, internal error) *AppError {
	return New(CodeServiceNotReady, http.StatusServiceUnavailable, message, internal)
}

func Internal(message string, internal error) *AppError {
	return New(CodeInternal, http.StatusInternalServerError, message, internal)
}
