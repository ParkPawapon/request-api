package errors

type Code string

const (
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeValidation      Code = "VALIDATION"
	CodeRateLimited     Code = "RATE_LIMITED"
	CodeTimeout         Code = "TIMEOUT"
	CodeInternal        Code = "INTERNAL"
	CodeServiceNotReady Code = "SERVICE_NOT_READY"
)
