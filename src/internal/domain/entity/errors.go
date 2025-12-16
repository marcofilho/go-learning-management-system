package entity

import "errors"

var (
	ErrNotFound                = errors.New("resource not found")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidInput            = errors.New("invalid input")
	ErrDuplicateEntry          = errors.New("duplicate entry")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrCourseNotAvailable      = errors.New("course not available")
	ErrAlreadyEnrolled         = errors.New("already enrolled in course")
)
