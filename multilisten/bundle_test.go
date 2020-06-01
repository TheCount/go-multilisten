package multilisten

import (
	"net"
	"testing"
)

// testListenAddr is a test listen address.
const testListenAddr = "localhost:47831"

// expectErr checks whether err has type Error.
func expectErr(t *testing.T, err error) {
	if err == nil {
		t.Fatal("Expected non-nil error")
	}
	if _, ok := err.(Error); !ok {
		t.Fatal("Expected error to be of type Error")
	}
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
