package metastorage

import (
	"context"
	"time"
)

// QueueState represents the state of a message in the queue
type QueueState int

const (
	StateIncoming QueueState = iota
	StateActive
	StateDeferred
	StateHold
	StateBounce
)

// String returns the string representation of the queue state
func (s QueueState) String() string {
	switch s {
	case StateIncoming:
		return "incoming"
	case StateActive:
		return "active"
	case StateDeferred:
		return "deferred"
	case StateHold:
		return "hold"
	case StateBounce:
		return "bounce"
	default:
		return "unknown"
	}
}

// MessageMetadata contains metadata about a message
type MessageMetadata struct {
	ID          string
	State       QueueState
	Attempts    int
	MaxAttempts int
	NextRetry   time.Time
	Created     time.Time
	Updated     time.Time
	LastError   string
	Size        int64
	Priority    int
	Headers     map[string]string
}

// MessageListOptions contains options for listing messages
type MessageListOptions struct {
	Limit     int       // Maximum number of messages to return
	Offset    int       // Number of messages to skip
	SortBy    string    // Sort field: "created", "updated", "priority", "attempts"
	SortOrder string    // Sort order: "asc" or "desc"
	Since     time.Time // Only return messages created/updated after this time
}

// MessageListResult contains the result of listing messages
type MessageListResult struct {
	MessageIDs []string // List of message IDs
	Total      int      // Total number of messages matching criteria
	HasMore    bool     // Whether there are more messages available
}

// Backend represents a metadata storage backend for message metadata
type Backend interface {
	// StoreMeta stores message metadata
	StoreMeta(ctx context.Context, messageID string, metadata MessageMetadata) error

	// GetMeta retrieves message metadata
	GetMeta(ctx context.Context, messageID string) (MessageMetadata, error)

	// UpdateMeta updates message metadata
	UpdateMeta(ctx context.Context, messageID string, metadata MessageMetadata) error

	// DeleteMeta removes message metadata
	DeleteMeta(ctx context.Context, messageID string) error

	// ListMessages lists messages with pagination and filtering
	ListMessages(ctx context.Context, state QueueState, options MessageListOptions) (MessageListResult, error)

	// NewMessageIterator creates an iterator for messages in a specific state
	// batchSize controls how many messages are fetched per backend call (optimization hint)
	NewMessageIterator(ctx context.Context, state QueueState, batchSize int) (MessageIterator, error)

	// MoveToState moves a message from one queue state to another atomically
	MoveToState(ctx context.Context, messageID string, fromState, toState QueueState) error

	// Close closes the metadata storage backend
	Close() error
}

// StateCounterBackend extends Backend with fast state counting capabilities
type StateCounterBackend interface {
	Backend
	
	// GetStateCount returns the cached count for a specific state (fast operation)
	GetStateCount(state QueueState) int64
}


// MessageIterator provides streaming access to messages in a specific state
type MessageIterator interface {
	// Next returns the next message metadata, whether more messages are available, and any error
	// Returns (metadata, hasMore, error)
	// When hasMore is false, the iterator is exhausted
	Next(ctx context.Context) (MessageMetadata, bool, error)
	
	// Close closes the iterator and releases any resources
	Close() error
}

// Factory creates metadata storage backends
type Factory interface {
	// Create creates a new metadata storage backend
	Create() (Backend, error)

	// Name returns the factory name
	Name() string
}
