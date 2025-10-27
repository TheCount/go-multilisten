package multilisten

import (
	"net"
)

// acceptInfo contains the result of an accept call.
type acceptInfo struct {
	// conn is the connection that was accepted. If the previous accept call
	// resulted in an error, conn is nil.
	conn net.Conn

	// err is the error the last Accept call returned. If there was no error,
	// err is nil.
	err *wrappedError

	// recovered is non-nil if the accepting goroutine exited due to a panic.
	// The goroutine will have recovered, so that the main Accept can re-raise
	// the panic.
	recovered any
}
