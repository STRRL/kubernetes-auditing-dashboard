package lifecycle

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrResourceNotFound indicates that no audit events were found for the specified resource
	ErrResourceNotFound = errors.New("resource not found in audit log")

	// ErrInvalidGVK indicates that the provided Group-Version-Kind format is invalid
	ErrInvalidGVK = errors.New("invalid GVK format")

	// ErrInvalidName indicates that the resource name is empty or invalid
	ErrInvalidName = errors.New("invalid resource name")

	// ErrDatabaseQuery indicates a database query error
	ErrDatabaseQuery = errors.New("database query failed")

	// ErrYAMLParsing indicates YAML parsing failure
	ErrYAMLParsing = errors.New("failed to parse YAML")

	// ErrJSONMarshaling indicates JSON marshaling failure
	ErrJSONMarshaling = errors.New("failed to marshal JSON")
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// DatabaseError wraps database-related errors with additional context
type DatabaseError struct {
	Operation string
	Query     string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation, query string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Query:     query,
		Err:       err,
	}
}

// ParseError wraps parsing-related errors
type ParseError struct {
	Type   string // "YAML", "JSON", "GVK", etc.
	Input  string
	Reason string
	Err    error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse %s: %s", e.Type, e.Reason)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// NewParseError creates a new parse error
func NewParseError(parseType, input, reason string, err error) *ParseError {
	return &ParseError{
		Type:   parseType,
		Input:  input,
		Reason: reason,
		Err:    err,
	}
}

// IsNotFound checks if an error indicates a resource was not found
func IsNotFound(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsDatabaseError checks if an error is a database error
func IsDatabaseError(err error) bool {
	var dbErr *DatabaseError
	return errors.As(err, &dbErr)
}

// IsParseError checks if an error is a parse error
func IsParseError(err error) bool {
	var parseErr *ParseError
	return errors.As(err, &parseErr)
}

// WrapError wraps an error with additional context
func WrapError(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", message, err)
}