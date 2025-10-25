package multilisten

import (
	"errors"
	"net"
	"strings"
	"testing"
)

// TestBasicErrorListener tests the listener field in a basicError.
func TestBasicErrorListener(t *testing.T) {
	l, err := net.Listen("tcp", testListenAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	bErr := &basicError{}
	if bErr.Listener() != nil {
		t.Error("Expected Listener() to return nil")
	}
	bErr.listener = l
	if bErr.Listener() != l {
		t.Error("Expected Listener() to return set listener")
	}
}

// TestBasicErrorStopped tests the stopped field in a basicError.
func TestBasicErrorStopped(t *testing.T) {
	bErr := &basicError{}
	if bErr.Stopped() {
		t.Error("Expected Stopped() to return false")
	}
	bErr.stopped = true
	if !bErr.Stopped() {
		t.Error("Expected Stopped() to return true")
	}
}

// TestGenericError tests generic errors.
func TestGenericError(t *testing.T) {
	const testMsg = "test error msg"
	var err net.Error
	err = &genericError{
		msg: testMsg,
	}
	if err.Error() != testMsg {
		t.Errorf("Bad error message: %s", err.Error())
	}
	if err.Timeout() {
		t.Error("generic error flagged as timeout")
	}
	if err.Temporary() {
		t.Error("generic error flagged as temporary")
	}
}

// TestWrappedError tests a wrapped error.
func TestWrappedError(t *testing.T) {
	const testMsg = "wrapped"
	var err1 net.Error
	err1 = &wrappedError{
		op:      "test",
		wrapped: errors.New(testMsg),
	}
	if !strings.HasSuffix(err1.Error(), testMsg) {
		t.Errorf("Bad error message: %s", err1.Error())
	}
	if err1.Timeout() {
		t.Errorf("Unexpected timeout error: %s", err1)
	}
	var err2 net.Error
	err2 = &wrappedError{
		op: "test",
		wrapped: &net.DNSError{
			Err:         testMsg,
			Name:        "test",
			Server:      "localhost",
			IsTimeout:   true,
			IsTemporary: false,
		},
		temporary: true,
	}
	if !err2.Timeout() {
		t.Error("Expected timeout error")
	}
	if !err2.Temporary() {
		t.Error("Expected temporary error")
	}
}
