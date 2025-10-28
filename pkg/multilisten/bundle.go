package multilisten

import (
	"net"
	"sync"
	"sync/atomic"
)

// bundle represents a bundle of listeners. It implements the net.Listener
// interface.
type bundle struct {
	// active is the number of currently active listeners.
	active atomic.Int64

	// once guards starting the Accept goroutines.
	once sync.Once

	// addr is the address of this listener.
	addr net.Addr

	// listeners is the bundle of listeners
	listeners []net.Listener

	// info is used by the accepting goroutines to report accepted connections
	// and errors.
	info chan acceptInfo
}

// runAccept accepts connections from the specified listener until an
// unrecoverable error occurs. Accepted connections or errors are communicated
// via the info channel.
func (b *bundle) runAccept(l net.Listener) {
	defer func() {
		if r := recover(); r != nil {
			b.info <- acceptInfo{
				recovered: r,
			}
		}
	}()
	for {
		var info acceptInfo
		var err error
		info.conn, err = l.Accept()
		if err != nil {
			info.err = &wrappedError{
				basicError: basicError{
					listener: l,
				},
				op:      "accept",
				wrapped: err,
			}
		}
		b.info <- info
		if info.err != nil {
			return
		}
	}
}

// start starts the accepting goroutines. This method should be called at most
// once.
func (b *bundle) start() {
	b.info = make(chan acceptInfo, len(b.listeners))
	b.active.Store(int64(len(b.listeners)))
	for _, l := range b.listeners {
		go b.runAccept(l)
	}
}

// Accept implements net.Listener.
func (b *bundle) Accept() (net.Conn, error) {
	b.once.Do(b.start)
	info, ok := <-b.info
	if !ok {
		return nil, ErrNoMoreListeners
	}
	if info.recovered != nil {
		if b.active.Add(-1) == 0 {
			close(b.info)
		}
		panic(info.recovered)
	}
	if info.err != nil {
		if b.active.Add(-1) == 0 {
			close(b.info)
		}
		return nil, info.err
	}

	return info.conn, nil
}

// Addr implements net.Listener.
func (b *bundle) Addr() net.Addr {
	return b.addr
}

// Close implements net.Listener.
func (b *bundle) Close() error {
	var err error
	for _, l := range b.listeners {
		if err2 := l.Close(); err2 != nil && err == nil {
			err = &wrappedError{
				basicError: basicError{
					listener: l,
				},
				op:      "close",
				wrapped: err2,
			}
		}
	}
	return err
}

// Bundle bundles multiple listeners into one. The specified main listener
// provides the address of the returned bundled listener.
//
// The Accept method of the returned listener will call Accept on all listeners
// concurrently and return the first connection thus obtained. If a sub-Accept
// call returns an error, listening on that listener will stop.
// The main accept method will wrap that the error
// into an [Error] and return it.
//
// The Close method of the returned listener will close all underlying
// listeners. Close returns the first error it encounters, or nil if none.
//
// Also consider using [NewSet] instead,
// which provides a more flexible API.
func Bundle(main net.Listener, more ...net.Listener) (net.Listener, error) {
	if main == nil {
		return nil, &genericError{
			msg: "main listener is nil",
		}
	}
	for _, l := range more {
		if l == nil {
			return nil, &genericError{
				msg: "listener is nil",
			}
		}
	}
	return &bundle{
		addr:      main.Addr(),
		listeners: append(more, main),
	}, nil
}
