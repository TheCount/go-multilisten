package multilisten

import (
	"net"
	"testing"
	"time"
)

// testListenAddr is a test listen address.
const testListenAddr = "localhost:47831"

// expectErr checks whether err has type Error.
// If so, it is returned.
func expectErr(t *testing.T, err error) Error {
	if err == nil {
		t.Fatal("Expected non-nil error")
	}
	x, ok := err.(Error)
	if !ok {
		t.Fatal("Expected error to be of type Error")
	}
	return x
}

// expectPermanentErr checks whether err is a permanent error of type Error.
func expectPermanentErr(t *testing.T, err error) {
	unpacked := expectErr(t, err)
	if unpacked.Temporary() {
		t.Fatal("Expected permanent error")
	}
}

// bundleSingleListener returns a bundle with a single listener.
func bundleSingleListener(t *testing.T) net.Listener {
	l, err := net.Listen("tcp", testListenAddr)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Bundle(l)
	if err != nil {
		l.Close()
		t.Fatalf("Bundling single listener failed: %s", err)
	}
	return b
}

// TestBundleNil tests calling bundle with nil listeners.
func TestBundleNil(t *testing.T) {
	_, err := Bundle(nil)
	expectErr(t, err)
	l, err := net.Listen("tcp", testListenAddr)
	if err != nil {
		t.Fatal(l)
	}
	defer l.Close()
	_, err = Bundle(l, nil)
	expectErr(t, err)
	_, err = Bundle(nil, l)
	expectErr(t, err)
}

// TestAddr tests the address function of the bundled listener.
func TestAddr(t *testing.T) {
	l, err := net.Listen("tcp", testListenAddr)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Bundle(l)
	if err != nil {
		l.Close()
		t.Fatalf("Bundling single listener failed: %s", err)
	}
	defer b.Close()
	origAddr := l.Addr()
	bAddr := b.Addr()
	if bAddr.Network() != origAddr.Network() {
		t.Errorf("Networks do not match: '%s' != '%s'",
			bAddr.Network(), origAddr.Network())
	}
	if bAddr.String() != origAddr.String() {
		t.Errorf("Addresses do not match: '%s' != '%s'", bAddr, origAddr)
	}
}

// TestCloseBeforeAccept tests what happens when the bundled listener is closed
// before Accept is called.
func TestCloseBeforeAccept(t *testing.T) {
	b := bundleSingleListener(t)
	if err := b.Close(); err != nil {
		t.Errorf("Closing bundled listener failed: %s", err)
	}
	_, err := b.Accept()
	expectPermanentErr(t, err)
}

// TestCloseWhileAccept tests calling Close while Accept is in progress.
func TestCloseWhileAccept(t *testing.T) {
	b := bundleSingleListener(t)
	done := make(chan error)
	go func() {
		_, err := b.Accept()
		done <- err
	}()
	time.Sleep(time.Second)
	if err := b.Close(); err != nil {
		t.Fatalf("Error closing bundled listener: %s", err)
	}
	err := <-done
	expectPermanentErr(t, err)
}
