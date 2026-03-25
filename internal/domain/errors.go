package domain

import "errors"

// User errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidRole       = errors.New("invalid role")
)

// Project type errors
var (
	ErrProjectTypeNotFound      = errors.New("project type not found")
	ErrProjectTypeAlreadyExists = errors.New("project type already exists")
	ErrInvalidProjectTypeName   = errors.New("project type name is between 5 and 100 characters")
)

// Application errors
var (
	ErrApplicationNotFound    = errors.New("application not found")
	ErrApplicationNotPending  = errors.New("application is not in pending status")
	ErrInvalidApplicationData = errors.New("invalid application data")
)

// Validation errors
var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrInvalidPhone = errors.New("invalid phone format")
	ErrInvalidURL   = errors.New("invalid URL format")
)

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)
