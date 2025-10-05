package net

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/extism/go-pdk"
)

var DefaultDialer = &Dialer{}

type Dialer struct{}

func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func (d *Dialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.DialContext(ctx, network, address)
}

//go:wasmimport wape:host/env net.dial
func _dial(network, addr uint64) int32

func (d *Dialer) DialContext(_ context.Context, network, addr string) (net.Conn, error) {
	networkMem := pdk.AllocateString(network)
	defer networkMem.Free()

	addrMem := pdk.AllocateString(addr)
	defer addrMem.Free()

	connID := _dial(networkMem.Offset(), addrMem.Offset())
	if connID < 0 {
		return nil, fmt.Errorf("failed to dial: %d", connID)
	}

	return &Conn{
		connID: connID,
	}, nil
}
