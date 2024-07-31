package lock

import "errors"

var (
	// ErrFailedToCreateLock is returned when the locker cannot be created
	ErrFailedToCreateLock = errors.New("failed to create locker")
	// ErrLockerClosed is returned when the locker is closed
	ErrLockerClosed = errors.New("locker is closed")
	// ErrNotLocked is returned when the lock is not acquired
	ErrNotLocked = errors.New("not locked")
)
