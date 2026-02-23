package metastorage

import (
	"errors"
)

// Common errors
var (
	// ErrMessageNotFound is returned when a message is not found
	ErrMessageNotFound = errors.New("message not found")

	// ErrBackendClosed is returned when the backend is closed
	ErrBackendClosed = errors.New("backend is closed")

	// ErrInvalidState is returned when an invalid state transition is requested
	ErrInvalidState = errors.New("invalid state transition")

	// ErrStateConflict is returned when a CAS operation fails because
	// the message is not in the expected fromState.
	// This is the expected outcome when multiple workers race for the same message.
	ErrStateConflict = errors.New("state conflict: message not in expected state")
)
