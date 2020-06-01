package multilisten

import (
	"net"
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
