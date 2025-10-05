package net

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"slices"

	extism "github.com/extism/go-sdk"

	"github.com/mymmrac/wape/internal"
)

// DialConfig configures [Dial].
type DialConfig struct {
	// NetworksAllowed configures the allowed network protocols. See [net.Dial] for allowed protocols. Defaults to none.
	NetworksAllowed []string
	// NetworksAllowAll allows all network protocols. Defaults to false.
	NetworksAllowAll bool

	// NetworkAddressesAllowed configures the allowed network addresses. Defaults to none.
	NetworkAddressesAllowed []string
	// NetworkAddressesAllowAll allows all network addresses. Defaults to false.
	NetworkAddressesAllowAll bool
}

// Dial creates a host function that calls [net.Dial].
func Dial(cfg DialConfig) extism.HostFunction {
	return internal.NewHostFunction("net.dial",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			network, err := p.ReadString(stack[0])
			if err != nil {
				panic(err)
			}

			if !cfg.NetworksAllowAll && !slices.Contains(cfg.NetworksAllowed, network) {
				panic(fmt.Errorf("network not allowed: %s", network))
			}

			addr, err := p.ReadString(stack[1])
			if err != nil {
				panic(err)
			}

			if !cfg.NetworkAddressesAllowAll && !slices.Contains(cfg.NetworkAddressesAllowed, addr) {
				panic(fmt.Errorf("address not allowed: %s", addr))
			}

			dialer := &net.Dialer{}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				panic(err)
			}

			connID := rand.Int32N(10000)
			internal.Connections.Set(connID, conn)

			stack[0] = extism.EncodeI32(connID)
		},
		[]extism.ValueType{extism.ValueTypePTR /* network */, extism.ValueTypePTR /* address */},
		[]extism.ValueType{extism.ValueTypeI32 /* connectionID | errorCode */},
	)
}
