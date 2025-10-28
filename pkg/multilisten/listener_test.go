package multilisten

import (
	"errors"
	"math/rand/v2"
	"net"
	"sync/atomic"
	"time"
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

func newListener() net.Listener {
	return &mockListener{
		addr: func() net.Addr { return mockAddr("ListenerMcEarface:2345") },
		accept: func() (net.Conn, error) {
			return mockConn{}, nil
		},
		close: func() error { return nil },
	}
}

// newPanickyListener returns a listener which panics when Accept is called.
func newPanickyListener() net.Listener {
	return &mockListener{
		addr:  func() net.Addr { return nil },
		close: func() error { return nil },
	}
}

// newCloseErrorListener returns a listener whose Close method returns the
// specified error.
func newCloseErrorListener(err error) net.Listener {
	return &mockListener{
		addr:   func() net.Addr { return nil },
		accept: func() (net.Conn, error) { return nil, nil },
		close:  func() error { return err },
	}
}

// newRandomListener creates a listener with random Accept behaviour. For each
// call to accept, the returned listener will increment numAccepts.
// Each Accept call takes a uniform random time from zero to maxAcceptTime.
// The failChance is the chance that the accept method returns an error.
func newRandomListener(
	numAccepts *atomic.Int32, maxAcceptTime time.Duration, failChance float32,
) net.Listener {
	done := false
	var closed atomic.Bool
	return &mockListener{
		addr: func() net.Addr { return nil },
		accept: func() (net.Conn, error) {
			defer numAccepts.Add(1)
			if closed.Load() {
				return nil, &genericError{msg: "listener closed"}
			}
			if done {
				panic("called again despite done being true")
			}
			time.Sleep(time.Duration(float32(maxAcceptTime) * rand.Float32()))
			if rand.Float32() < failChance {
				err := &wrappedError{
					op:      "accept",
					wrapped: errors.New("failed"),
				}
				done = true
				return nil, err
			}
			return mockConn{}, nil
		},
		close: func() error {
			if closed.Load() {
				return &genericError{msg: "listener already closed"}
			}
			closed.Store(true)
			return nil
		},
	}
}
