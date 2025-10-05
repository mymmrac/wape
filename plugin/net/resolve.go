package net

import (
	"context"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

var DefaultResolver = &Resolver{}

type Resolver struct{}

//go:wasmimport wape:host/env net.resolver.lookupHost
func _lookupHost(host uint64) uint64

func (r *Resolver) LookupHost(_ context.Context, host string) ([]string, error) {
	hostMem := pdk.AllocateString(host)
	defer hostMem.Free()

	addressesPtr := _lookupHost(hostMem.Offset())
	if addressesPtr < 0 {
		return nil, fmt.Errorf("failed to lookup host: %d", addressesPtr)
	}

	addresses := pdk.ParamString(addressesPtr)
	return strings.Split(addresses, ","), nil
}
