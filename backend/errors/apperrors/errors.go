package apperrors

import "errors"

// ErrUserNotFound is returned when a user is not found in the database.
var ErrUserNotFound = errors.New("user not found")

// ErrInvalidCredentials is returned when the user provides invalid credentials.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrUserAlreadyExists is return when a user already exists in the database.
var ErrUserAlreadyExists = errors.New("user alraedy exists")

var ErrInvalidEnv = errors.New("invalid environment")
