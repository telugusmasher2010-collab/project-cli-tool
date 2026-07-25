package errors

import (
	"fmt"
)

type Code string

const (
	ErrTemplateNotFound Code = "TEMPLATE_NOT_FOUND"
	ErrGenerationFailed Code = "GENERATION_FAILED"
	ErrVariableMissing  Code = "VARIABLE_MISSING"
	ErrInvalidInput     Code = "INVALID_INPUT"
	ErrHookFailed       Code = "HOOK_FAILED"
	ErrPermissionDenied Code = "PERMISSION_DENIED"
	ErrInternal         Code = "INTERNAL_ERROR"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

func IsTemplateNotFound(err error) bool {
	var e *Error
	if ok := as(err, &e); ok {
		return e.Code == ErrTemplateNotFound
	}
	return false
}

func IsGenerationFailed(err error) bool {
	var e *Error
	if ok := as(err, &e); ok {
		return e.Code == ErrGenerationFailed
	}
	return false
}

func IsInvalidInput(err error) bool {
	var e *Error
	if ok := as(err, &e); ok {
		return e.Code == ErrInvalidInput
	}
	return false
}

func UserMessage(err error) string {
	var e *Error
	if ok := as(err, &e); ok {
		switch e.Code {
		case ErrTemplateNotFound:
			return fmt.Sprintf("Template not found: %s", e.Message)
		case ErrGenerationFailed:
			return fmt.Sprintf("Failed to generate project: %s", e.Message)
		case ErrVariableMissing:
			return fmt.Sprintf("Missing required variable: %s", e.Message)
		case ErrInvalidInput:
			return fmt.Sprintf("Invalid input: %s", e.Message)
		case ErrHookFailed:
			return fmt.Sprintf("Post-generation hook failed: %s", e.Message)
		case ErrPermissionDenied:
			return fmt.Sprintf("Permission denied: %s", e.Message)
		default:
			return fmt.Sprintf("Something went wrong: %s", e.Message)
		}
	}
	return err.Error()
}

func as(err error, target interface{}) bool {
	for {
		if e, ok := err.(*Error); ok {
			if t, ok2 := target.(**Error); ok2 {
				*t = e
				return true
			}
			return false
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
}
