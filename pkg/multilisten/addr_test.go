package multilisten

import (
	"net"
)

const mockNetwork = "tcp"

type mockAddr string

var _ net.Addr = mockAddr("")

func (mockAddr) Network() string {
	return mockNetwork
}

func (a mockAddr) String() string {
	return string(a)
}
