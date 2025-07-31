package metastorage

import (
	"errors"
)


// Common errors
var (
	ErrMessageNotFound   = errors.New("message not found")
	ErrInvalidState      = errors.New("invalid state transition")
	ErrBackendClosed     = errors.New("backend is closed")
)