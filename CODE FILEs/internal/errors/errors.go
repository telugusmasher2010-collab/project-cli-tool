// Package errors provides custom error types and error codes for the proj-init CLI.
//
// It defines a structured Error type that carries a machine-readable Code,
// a human-readable Message, and an optional wrapped Cause. This enables
// callers to programmatically inspect errors using IsCode or the standard
// errors.Is/errors.As functions.
package errors

import (
	"fmt"
)

// Code represents a machine-readable error category.
type Code string

const (
	// ErrTemplateNotFound indicates the requested template does not exist.
	ErrTemplateNotFound Code = "TEMPLATE_NOT_FOUND"
	// ErrTemplateExists indicates a template with the given name already exists.
	ErrTemplateExists Code = "TEMPLATE_EXISTS"
	// ErrInvalidTemplate indicates the template content is malformed or invalid.
	ErrInvalidTemplate Code = "INVALID_TEMPLATE"
	// ErrInvalidInput indicates the user provided invalid input.
	ErrInvalidInput Code = "INVALID_INPUT"
	// ErrConfigNotFound indicates the configuration file was not found.
	ErrConfigNotFound Code = "CONFIG_NOT_FOUND"
	// ErrConfigInvalid indicates the configuration file is malformed or invalid.
	ErrConfigInvalid Code = "CONFIG_INVALID"
	// ErrGenerationFailed indicates project generation failed.
	ErrGenerationFailed Code = "GENERATION_FAILED"
	// ErrFilesystem indicates a filesystem operation failed.
	ErrFilesystem Code = "FILESYSTEM_ERROR"
	// ErrInternal indicates an unexpected internal error.
	ErrInternal Code = "INTERNAL_ERROR"
	// ErrOutputExists indicates the output directory already exists.
	ErrOutputExists Code = "OUTPUT_EXISTS"
)

// Error is a structured error that carries a Code, Message, and optional Cause.
// It implements the error interface and supports unwrapping via Unwrap.
type Error struct {
	Code    Code
	Message string
	Err     error
}

// Error returns a formatted error string. If a wrapped cause exists it is
// included in the output.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause, enabling errors.Is and errors.As
// to traverse the error chain.
func (e *Error) Unwrap() error {
	return e.Err
}

// New creates a new Error with the given code and message and no wrapped cause.
func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap creates a new Error with the given code, message, and an underlying
// cause that is preserved via Unwrap for use with errors.Is/errors.As.
func Wrap(code Code, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// IsCode reports whether err (or any error in its chain) is an *Error
// whose Code matches the given code.
func IsCode(err error, code Code) bool {
	for err != nil {
		if e, ok := err.(*Error); ok {
			if e.Code == code {
				return true
			}
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// IsTemplateNotFound reports whether err indicates a missing template.
func IsTemplateNotFound(err error) bool {
	return IsCode(err, ErrTemplateNotFound)
}

// IsGenerationFailed reports whether err indicates a generation failure.
func IsGenerationFailed(err error) bool {
	return IsCode(err, ErrGenerationFailed)
}

// IsInvalidInput reports whether err indicates invalid user input.
func IsInvalidInput(err error) bool {
	return IsCode(err, ErrInvalidInput)
}

// IsOutputExists reports whether err indicates the output directory already exists.
func IsOutputExists(err error) bool {
	return IsCode(err, ErrOutputExists)
}

// IsFilesystem reports whether err indicates a filesystem operation failure.
func IsFilesystem(err error) bool {
	return IsCode(err, ErrFilesystem)
}

// UserMessage extracts a user-friendly message from a structured Error.
// If err is not a structured Error, its default Error() string is returned.
func UserMessage(err error) string {
	if err == nil {
		return ""
	}
	for {
		if e, ok := err.(*Error); ok {
			switch e.Code {
			case ErrTemplateNotFound:
				return fmt.Sprintf("Template not found: %s", e.Message)
			case ErrTemplateExists:
				return fmt.Sprintf("Template already exists: %s", e.Message)
			case ErrInvalidTemplate:
				return fmt.Sprintf("Invalid template: %s", e.Message)
			case ErrGenerationFailed:
				return fmt.Sprintf("Failed to generate project: %s", e.Message)
			case ErrInvalidInput:
				return fmt.Sprintf("Invalid input: %s", e.Message)
			case ErrConfigNotFound:
				return fmt.Sprintf("Config not found: %s", e.Message)
			case ErrConfigInvalid:
				return fmt.Sprintf("Invalid config: %s", e.Message)
		case ErrFilesystem:
			return fmt.Sprintf("Filesystem error: %s", e.Message)
		case ErrOutputExists:
			return fmt.Sprintf("Output already exists: %s", e.Message)
			default:
				return fmt.Sprintf("Something went wrong: %s", e.Message)
			}
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return err.Error()
}
