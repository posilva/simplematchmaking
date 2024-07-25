package services

import "fmt"

var (
	// ErrMatchNotFound is the error returned when the match is not found
	ErrMatchNotFound = fmt.Errorf("match not found")
)
