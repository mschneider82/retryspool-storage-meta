# RetrySpool Meta Storage

Metadata storage backend interfaces for the RetrySpool message queue system. This package handles the storage and retrieval of message metadata.

## Overview

This package provides interfaces for storing and querying message metadata efficiently. It's designed to work with backends like etcd, Redis, PostgreSQL, etc.

## Installation

```bash
go get schneider.vip/retryspool/storage/meta
```

## Interfaces

### Backend

The core metadata storage interface for message metadata:

```go
type Backend interface {
    StoreMeta(ctx context.Context, messageID string, metadata MessageMetadata) error
    GetMeta(ctx context.Context, messageID string) (MessageMetadata, error)
    UpdateMeta(ctx context.Context, messageID string, metadata MessageMetadata) error
    DeleteMeta(ctx context.Context, messageID string) error
    ListMessages(ctx context.Context, state QueueState, options MessageListOptions) (MessageListResult, error)
    MoveToState(ctx context.Context, messageID string, newState QueueState) error
    Close() error
}
```

### Extended Interfaces

For backends with additional capabilities:

```go
// Fast state counting for performance optimization
type StateCounterBackend interface {
    Backend
    GetStateCount(state QueueState) int64
}
```

### Factory

Factory pattern for creating metadata storage backends:

```go
type Factory interface {
    Create() (Backend, error)
    Name() string
}
```

## Usage

### Basic Usage

```go
import metastorage "schneider.vip/retryspool/storage/meta"

// Create a factory (implementation-specific)
factory := filesystem.NewFactory("/path/to/metadata")

// Create backend
backend, err := factory.Create()
if err != nil {
    panic(err)
}
defer backend.Close()

// Store metadata
metadata := metastorage.MessageMetadata{
    ID:          "msg-123",
    State:       metastorage.StateIncoming,
    Attempts:    0,
    MaxAttempts: 3,
    Created:     time.Now(),
    Headers:     map[string]string{"to": "user@example.com"},
}

err = backend.StoreMeta(ctx, "msg-123", metadata)
if err != nil {
    panic(err)
}

// Query metadata
meta, err := backend.GetMeta(ctx, "msg-123")
if err != nil {
    panic(err)
}

// List messages
result, err := backend.ListMessages(ctx, metastorage.StateIncoming, metastorage.MessageListOptions{
    Limit:  10,
    Offset: 0,
    SortBy: "created",
})
```

### State Counter Usage

```go
// Check if backend supports fast state counting
if stateCounter, ok := backend.(metastorage.StateCounterBackend); ok {
    count := stateCounter.GetStateCount(metastorage.StateIncoming)
    fmt.Printf("Messages in incoming state: %d\n", count)
} else {
    // Fallback to listing messages for count
    result, err := backend.ListMessages(ctx, metastorage.StateIncoming, metastorage.MessageListOptions{
        Limit: 0, // Get total count only
    })
    if err == nil {
        fmt.Printf("Messages in incoming state: %d\n", result.Total)
    }
}
```


## Design Principles

- **Separation of Concerns**: Only handles message metadata, not data
- **Query Optimized**: Efficient listing, counting, and filtering
- **Backend Agnostic**: Works with any metadata storage backend
- **Atomic Operations**: Built-in atomic state transitions and updates
- **State Management**: Built-in queue state transitions
- **Performance**: Optimized for frequent metadata queries

## Available Implementations

- **Filesystem**: `schneider.vip/retryspool/storage/meta/filesystem`
- **etcd**: (planned)
- **Redis**: (planned)
- **PostgreSQL**: (planned)
- **SQLite**: (planned)

## Performance Considerations

- Use appropriate indexing for your backend (state, created, priority)
- Consider connection pooling for database backends
- Use atomic operations for state transitions to ensure consistency
- Implement proper pagination for large message lists
- Cache frequently accessed metadata if needed

## Queue States

The system supports these queue states:

- `StateIncoming` - New messages awaiting processing
- `StateActive` - Messages currently being processed  
- `StateDeferred` - Messages waiting for retry
- `StateHold` - Messages manually held (admin intervention)
- `StateBounce` - Messages that permanently failed