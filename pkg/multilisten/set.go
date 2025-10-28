package multilisten

import (
	"net"
)

// Set is a set of listeners that can be accepted from as a single listener.
//
// The methods of Set are safe for concurrent use by multiple goroutines.
type Set struct {
	closed     chan struct{}
	acceptChan chan acceptInfo
	addr       net.Addr
}

var _ net.Listener = (*Set)(nil)

// Accept waits for and returns the next connection
// to any of the listeners in the set.
func (s *Set) Accept() (net.Conn, error) {
	select {
	case info := <-s.acceptChan:
		if info.recovered != nil {
			panic(info.recovered)
		}
		if info.err != nil {
			return nil, info.err
		}
		return info.conn, nil
	case <-s.closed:
		return nil, net.ErrClosed
	}
}

// Addr returns the address set with [NewSet].
func (s *Set) Addr() net.Addr {
	return s.addr
}

// Close closes this set of listeners.
//
// Note that while Close will cause all underlying listeners to be closed
// (as this is the only way to stop Accept calls), it does not return
// potential errors from closing those listeners.
func (s *Set) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = net.ErrClosed
		}
	}()

	close(s.closed)
	return nil
}

// Add adds a listener to the set.
func (s *Set) Add(l net.Listener) error {
	if l == nil {
		panic("nil listener")
	}

	select {
	case <-s.closed:
		return net.ErrClosed
	default:
	}

	stopped := make(chan struct{})

	go func() {
		select {
		case <-s.closed:
			l.Close()
		case <-stopped:
		}
	}()

	go func() {
		defer close(stopped)

		defer func() {
			if r := recover(); r != nil {
				select {
				case s.acceptChan <- acceptInfo{
					recovered: r,
				}:
				case <-s.closed:
				}
			}
		}()

		for {
			var acceptInfo acceptInfo
			if conn, err := l.Accept(); err != nil {
				acceptInfo.err = &wrappedError{
					basicError: basicError{
						listener: l,
					},
					op:      "accept",
					wrapped: err,
				}
			} else {
				acceptInfo.conn = conn
			}

			select {
			case s.acceptChan <- acceptInfo:
			case <-s.closed:
				if conn := acceptInfo.conn; conn != nil {
					// edge case: accept succeeded just as Set was closed.
					conn.Close()
				}
				return
			}
		}
	}()

	return nil
}

// InjectConn causes a subsequent call to [Accept] to return the given
// connection.
func (s *Set) InjectConn(conn net.Conn) error {
	if conn == nil {
		panic("nil conn")
	}

	select {
	case s.acceptChan <- acceptInfo{conn: conn}:
		return nil
	default:
	}

	select {
	case <-s.closed:
		return net.ErrClosed
	default:
	}

	go func() {
		select {
		case s.acceptChan <- acceptInfo{conn: conn}:
		case <-s.closed:
			conn.Close()
		}
	}()

	return nil
}

// NewSet creates a new Set, empty set of listeners, with the given address.
func NewSet(addr net.Addr) *Set {
	if addr == nil {
		panic("nil addr")
	}

	return &Set{
		closed:     make(chan struct{}),
		acceptChan: make(chan acceptInfo),
		addr:       addr,
	}
}
