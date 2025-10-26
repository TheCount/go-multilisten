package multilisten

import (
	"fmt"
	"net"
)

// ErrNoMoreListeners is returned when all listeners in a bundle have stopped.
var ErrNoMoreListeners Error = &genericError{
	msg: "all listeners stopped",
}

// Error is the interface implemented by errors produced by this package.
type Error interface {
	net.Error

	// Listener returns the listener on which the error originally occurred.
	// If the error is not listener-specific, this is nil.
	Listener() net.Listener
}

// basicError is basic component of all multilisten errors.
type basicError struct {
	// listener is the listener which caused the error, or nil.
	listener net.Listener
}

// Listener implements Error.
func (err *basicError) Listener() net.Listener {
	return err.listener
}

// genericError is the generic error structure for non-temporary, non-timeout
// listener errors not caused by underlying errors.
type genericError struct {
	basicError

	// msg is the error message.
	msg string
}

// Error returns the error message.
func (err *genericError) Error() string {
	return err.msg
}

// Timeout always returns false.
func (err *genericError) Timeout() bool {
	return false
}

// Temporary always returns false.
func (err *genericError) Temporary() bool {
	return false
}

// wrappedError wraps another error.
type wrappedError struct {
	basicError

	// op is the operation during which the error occurred.
	op string

	// wrapped is the wrapped error.
	wrapped error
}

// Error implements error.
func (err *wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", err.op, err.wrapped.Error())
}

// Timeout returns true if and only if the wrapped error is a timeout error.
func (err *wrappedError) Timeout() bool {
	return isTimeout(err.wrapped)
}

// Temporary returns whether this is a temporary error.
func (err *wrappedError) Temporary() bool {
	// Note: as of Go 1.18, net.Error.Temporary is deprecated.
	return err.Timeout()
}

// Unwrap returns the wrapped error.
func (err *wrappedError) Unwrap() error {
	return err.wrapped
}

// isTimeout returns true if and only if the specified error has a Timeout
// method which returns true.
func isTimeout(err error) bool {
	if x, ok := err.(interface{ Timeout() bool }); ok {
		return x.Timeout()
	}
	return false
}
