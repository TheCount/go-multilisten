package multilisten

import (
	"errors"
	"math/rand"
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
// call to accept, the returned listener will increment *numAccepts atomically.
// Each Accept call takes a uniform random time from zero to maxAcceptTime.
// The failChance is the chance that the accept method returns an error.
// On average half of these errors will be permanent.
func newRandomListener(
	numAccepts *int32, maxAcceptTime time.Duration, failChance float32,
) net.Listener {
	done := false
	var closed int32
	return &mockListener{
		addr: func() net.Addr { return nil },
		accept: func() (net.Conn, error) {
			defer atomic.AddInt32(numAccepts, 1)
			if atomic.LoadInt32(&closed) != 0 {
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
				if rand.Float32() < 0.5 {
					err.temporary = true
				} else {
					done = true
				}
				return nil, err
			}
			return mockConn{}, nil
		},
		close: func() error {
			if atomic.LoadInt32(&closed) != 0 {
				return &genericError{msg: "listener already closed"}
			}
			atomic.StoreInt32(&closed, 1)
			return nil
		},
	}
}
