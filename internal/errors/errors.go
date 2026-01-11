package errors

import (
	"errors"
	"fmt"
)

// JRXError represents a custom error with operation context
type JRXError struct {
	Op  string // Operation that failed
	Err error  // Original error
}

func (e *JRXError) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *JRXError) Unwrap() error {
	return e.Err
}

// NewError creates a new JRXError
func NewError(op string, err error) *JRXError {
	return &JRXError{
		Op:  op,
		Err: err,
	}
}

// Common errors
var (
	ErrTemplateNotFound    = errors.New("template not found")
	ErrProjectExists       = errors.New("project directory already exists")
	ErrInvalidTemplate     = errors.New("invalid template configuration")
	ErrTemplatePathMissing = errors.New("template path does not exist")
	ErrConfigNotFound      = errors.New("configuration file not found")
	ErrEmptyProjectName    = errors.New("project name cannot be empty")
	ErrEmptyTemplateName   = errors.New("template name cannot be empty")
	ErrCloneRepository     = errors.New("failed to clone template repository")
	ErrLoadTemplates       = errors.New("failed to load templates")
)
