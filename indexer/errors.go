package indexer

import (
	"fmt"
)

// Error types for better error handling
type ErrorType int

const (
	ErrInvalidPath ErrorType = iota
	ErrInvalidXML
	ErrIndexNotFound
	ErrInvalidTerm
	ErrIOError
)

type WikiError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *WikiError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *WikiError) Unwrap() error {
	return e.Cause
}

// Error constructors
func NewInvalidPathError(path string, cause error) *WikiError {
	return &WikiError{
		Type:    ErrInvalidPath,
		Message: fmt.Sprintf("invalid path: %s", path),
		Cause:   cause,
	}
}

func NewInvalidXMLError(cause error) *WikiError {
	return &WikiError{
		Type:    ErrInvalidXML,
		Message: "invalid XML format",
		Cause:   cause,
	}
}

func NewIndexNotFoundError(path string) *WikiError {
	return &WikiError{
		Type:    ErrIndexNotFound,
		Message: fmt.Sprintf("index not found: %s", path),
	}
}

func NewInvalidTermError(term string) *WikiError {
	return &WikiError{
		Type:    ErrInvalidTerm,
		Message: fmt.Sprintf("invalid term: %s", term),
	}
}

func NewIOError(operation string, cause error) *WikiError {
	return &WikiError{
		Type:    ErrIOError,
		Message: fmt.Sprintf("IO error during %s", operation),
		Cause:   cause,
	}
}
