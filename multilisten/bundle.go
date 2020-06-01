package multilisten

import (
	"net"
	"sync"
	"sync/atomic"
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
	recovered interface{}
}

// bundle represents a bundle of listeners. It implements the net.Listener
// interface.
type bundle struct {
	// active is the number of currently active listeners.
	active int64

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
				op:        "accept",
				wrapped:   err,
				temporary: true, // assume true due to other possible listeners
			}
			if x, ok := err.(interface{ Temporary() bool }); !ok || !x.Temporary() {
				info.err.stopped = true
			}
		}
		b.info <- info
		if info.err != nil && info.err.stopped {
			return
		}
	}
}

// start starts the accepting goroutines. This method should be called at most
// once.
func (b *bundle) start() {
	b.info = make(chan acceptInfo, len(b.listeners))
	atomic.StoreInt64(&b.active, int64(len(b.listeners)))
	for _, l := range b.listeners {
		go b.runAccept(l)
	}
}

// Accept implements net.Listener.
func (b *bundle) Accept() (net.Conn, error) {
	b.once.Do(b.start)
	info, ok := <-b.info
	if !ok {
		return nil, &genericError{
			msg: "all listeners stopped",
		}
	}
	if info.recovered != nil {
		if atomic.AddInt64(&b.active, -1) == 0 {
			close(b.info)
		}
		panic(info.recovered)
	}
	if info.err != nil && info.err.stopped {
		if atomic.AddInt64(&b.active, -1) == 0 {
			close(b.info)
			info.err.temporary = false
		}
	}
	return info.conn, info.err
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
					stopped:  true,
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
// concurrently and return the first connection thus obtained. If a sub-Access
// call returns an error, listening on that listener will stop if and only
// if the error is not temporary. The main accept method will wrap the error
// into an Error and set it to non-temporary only if the sub-Accept was on the
// last remaining listener.
//
// The Close method of the returned listener will close all underlying
// listeners. Close returns the first error it encounters, or nil if none.
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
