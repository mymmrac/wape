package net

import (
	"context"
	"fmt"
	"net"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user net.dial
func _dial(network, addr uint64) int32

var DefaultDialer = &Dialer{}

type Dialer struct{}

func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func (d *Dialer) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	networkMem := pdk.AllocateString(network)
	defer networkMem.Free()

	addrMem := pdk.AllocateString(addr)
	defer addrMem.Free()

	connectionID := _dial(networkMem.Offset(), addrMem.Offset())
	if connectionID < 0 {
		return nil, fmt.Errorf("failed to dial: %d", connectionID)
	}

	return &Conn{
		connectionID: connectionID,
	}, nil
}
