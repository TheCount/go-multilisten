package multilisten

import (
	"net"
)

// mockListener mocks a listener.
type mockListener struct {
	addr   func() net.Addr
	accept func() (net.Conn, error)
	close  func() error
}

// Addr implements net.Listener.
func (l *mockListener) Addr() net.Addr {
	if l.addr == nil {
		panic("Addr called")
	}
	return l.addr()
}

// Accept implements net.Listener.
func (l *mockListener) Accept() (net.Conn, error) {
	if l.accept == nil {
		panic("Accept called")
	}
	return l.accept()
}

// Close implements net.Listener.
func (l *mockListener) Close() error {
	if l.close == nil {
		panic("Close called")
	}
	return l.close()
}
