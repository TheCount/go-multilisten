package multilisten

import (
	"errors"
	"net"
	"testing"
)

func TestNewSetNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for NewSet(nil)")
		}
	}()

	NewSet(nil)
}

func TestNewSet(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))
	if set == nil {
		t.Fatal("Expected non-nil NewSet return value")
	}
	if err := set.Close(); err != nil {
		t.Fatalf("Close new Set: %s", err)
	}
}

func TestSetAddr(t *testing.T) {
	addr := mockAddr("MockyMcMockface:1234")
	set := NewSet(addr)
	defer set.Close()

	if setAddr := set.Addr(); setAddr != addr {
		t.Fatalf("Expected Addr to be %q, got %q", addr, setAddr)
	}
}

func TestSetClose(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))

	if err := set.Close(); err != nil {
		t.Fatalf("Close new Set: %s", err)
	}

	if err := set.Close(); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("Expected second Close to return net.ErrClosed, got: %s", err)
	}
}

func TestSetAddNil(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))
	defer set.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for Set.Add(nil)")
		}
	}()

	set.Add(nil) //nolint panic expected
}

func TestSetAdd(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))
	defer set.Close()

	listener := newListener()
	if err := set.Add(listener); err != nil {
		t.Fatalf("Set.Add: %s", err)
	}
	if _, err := set.Accept(); err != nil {
		t.Fatalf("Set.Accept after Add: %s", err)
	}
}

func TestSetInjectNil(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))
	defer set.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for Set.InjectConn(nil)")
		}
	}()

	set.InjectConn(nil) //nolint panic expected
}

func TestSetInject(t *testing.T) {
	set := NewSet(mockAddr("MockyMcMockface:1234"))
	defer set.Close()

	conn := mockConn{}
	if err := set.InjectConn(&conn); err != nil {
		t.Fatalf("Set.InjectConn: %s", err)
	}
	accepted, err := set.Accept()
	if err != nil {
		t.Fatalf("Set.Accept after InjectConn: %s", err)
	}
	if accepted != &conn {
		t.Fatalf("Expected injected conn %p, got %p", &conn, accepted)
	}
}
