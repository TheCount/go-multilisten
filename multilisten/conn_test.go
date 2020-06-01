package multilisten

import (
	"net"
	"time"
)

// mockConn mocks a connection.
type mockConn struct{}

// Close implements net.Conn.
func (mockConn) Close() error {
	return nil
}

// LocalAddr implements net.Conn.
func (mockConn) LocalAddr() net.Addr {
	return nil
}

// Read implements net.Conn.
func (mockConn) Read(buf []byte) (int, error) {
	return len(buf), nil
}

// RemoteAddr implements net.Conn.
func (mockConn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline implements net.Conn.
func (mockConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn.
func (mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements net.Conn.
func (mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// Write implements net.Conn.
func (mockConn) Write(buf []byte) (int, error) {
	return len(buf), nil
}
