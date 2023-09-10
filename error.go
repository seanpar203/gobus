package gobus

import "fmt"

type InvalidArgError struct {
	Func     string
	Expected any
	Actual   any
}

// NewInvalidArgError creates a new InvalidArgError instance with the given parameters.
//
// Parameters:
// - event: the event associated with the error.
// - funcName: the name of the function where the error occurred.
// - expected: the expected value.
// - actual: the actual value.
//
// Returns:
// - InvalidArgError: a new instance of InvalidArgError.
func NewInvalidArgError(event Event, funcName string, expected any, actual any) *InvalidArgError {
	return &InvalidArgError{
		Func:     funcName,
		Expected: expected,
		Actual:   actual,
	}
}

// Error returns the error message for the InvalidArgError type.
//
// It formats the error message with the event, function, expected, and actual values.
// Returns a string.
func (e *InvalidArgError) Error() string {
	return fmt.Sprintf("Invalid argument for event func: %s. Expected %v, got %v", e.Func, e.Expected, e.Actual)
}
