package queues

import "errors"

var (
	// ErrFailedToAcquireLock is returned when the lock cannot be acquired
	ErrFailedToAcquireLock = errors.New("failed to acquire lock")

	// ErrFailedToParseValue is returned when the redis value cannot be parsed
	ErrFailedToParseValue = errors.New("failed to parse value")

	// ErrFailedToEncodeQueueEntry is returned when the queue entry cannot be marshaled
	ErrFailedToEncodeQueueEntry = errors.New("failed to encode queue entry")

	// ErrFailedToDecodeQueueEntry is returned when the queue entry cannot be marshaled
	ErrFailedToDecodeQueueEntry = errors.New("failed to decode queue entry")

	// ErrFailedExecuteCommand is returned when the redis command cannot be executed
	ErrFailedExecuteCommand = errors.New("failed to execute command")
)
