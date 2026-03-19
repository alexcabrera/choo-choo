package model

type ErrorCode string

const (
	ErrCrushFailed   ErrorCode = "crush_failed"
	ErrValidation    ErrorCode = "validation"
	ErrFileCorrupted ErrorCode = "file_corrupted"
	ErrMaxRetries    ErrorCode = "max_retries"
	ErrProcessCrash  ErrorCode = "process_crash"
)

type AppError struct {
	Code      ErrorCode
	Message   string
	Recovery  string
}

func (e AppError) Error() string {
	return string(e.Code) + ": " + e.Message
}

func NewAppError(code ErrorCode, message, recovery string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Recovery: recovery,
	}
}
